#!/bin/bash

# exit immediately if any command fails
set -e

# CONFIG
PROJECT_ID="project-192493f2-913e-4971-86c"
REGION="europe-west1"
ZONE="europe-west1-b"
CLUSTER="kuar-cluster"
REGISTRY="europe-west1-docker.pkg.dev"
REPO="inference-gateway"
IMAGE="gateway"
NAMESPACE="inference-gateway"
DEPLOYMENT="gateway"

# Version
VERSION=$(git rev-parse --short HEAD)
FULL_IMAGE="$REGISTRY/$PROJECT_ID/$REPO/$IMAGE:$VERSION"


log() {
  echo ""
  echo "→ $1"
}

error() {
  echo ""
  echo "✗ ERROR: $1"
  exit 1
}

success() {
  echo ""
  echo "✓ $1"
}

# ── Pre-flight checks ──────────────────────────────────
log "Running pre-flight checks..."

# Check required tools are installed
command -v docker   >/dev/null 2>&1 || error "docker not found"
command -v gcloud   >/dev/null 2>&1 || error "gcloud not found"
command -v kubectl  >/dev/null 2>&1 || error "kubectl not found"

# Check we're on the right git branch
BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$BRANCH" != "main" ]; then
  echo "  Warning: deploying from branch '$BRANCH', not main"
  read -p "  Continue? (y/n): " confirm
  if [ "$confirm" != "y" ]; then
    error "Deploy cancelled"
  fi
fi

# Check there are no uncommitted changes
if ! git diff --quiet; then
  echo "  Warning: you have uncommitted changes"
  read -p "  Continue anyway? (y/n): " confirm
  if [ "$confirm" != "y" ]; then
    error "Deploy cancelled — commit your changes first"
  fi
fi

success "Pre-flight checks passed"

log "Building Docker image..."
echo "  Image: $FULL_IMAGE"

docker build \
  --platform linux/amd64 \
  -t "$FULL_IMAGE" \
  .

success "Image built"

log "Pushing to Artifact Registry..."

gcloud auth configure-docker "$REGISTRY" --quiet

docker push "$FULL_IMAGE"

success "Image pushed"

# ── Deploy ────────────────────────────────────────────
log "Connecting to GKE cluster..."

gcloud container clusters get-credentials "$CLUSTER" \
  --zone "$ZONE" \
  --project "$PROJECT_ID" \
  --quiet

log "Deploying to GKE..."
echo "  Version: $VERSION"
echo "  Namespace: $NAMESPACE"

# Update the deployment image
# This triggers a rolling update automatically
kubectl set image deployment/"$DEPLOYMENT" \
  "$IMAGE=$FULL_IMAGE" \
  -n "$NAMESPACE"

# ── Wait for rollout ───────────────────────────────────
log "Waiting for rollout to complete..."

kubectl rollout status deployment/"$DEPLOYMENT" \
  -n "$NAMESPACE" \
  --timeout=120s

# ── Verify ────────────────────────────────────────────
log "Verifying deployment..."

kubectl get pods -n "$NAMESPACE"

# Get external IP
GATEWAY_IP=$(kubectl get service gateway \
  -n "$NAMESPACE" \
  -o jsonpath='{.status.loadBalancer.ingress[0].ip}')

# Hit health endpoint to confirm
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" \
  http://"$GATEWAY_IP"/health)

if [ "$HTTP_STATUS" = "200" ]; then
  success "Health check passed — gateway is live"
else
  error "Health check failed — got HTTP $HTTP_STATUS"
fi

# ── Summary ───────────────────────────────────────────
echo ""
echo "╔════════════════════════════════════════╗"
echo "║         Deployment Complete            ║"
echo "╚════════════════════════════════════════╝"
echo "→ Version:     $VERSION"
echo "→ Image:       $FULL_IMAGE"
echo "→ Gateway URL: http://$GATEWAY_IP"
echo "→ Health:      http://$GATEWAY_IP/health"
echo ""