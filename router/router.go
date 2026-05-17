package router

import (
	"inference-gateway/interfaces"
	"strings"
)

// Defines a single routing decision
type Rule struct {
	Name      string
	Condition func(prompt string) bool
	Provider  interfaces.ModelProvider
	Model     string
}

// Holds all routing rules and fallback provider
type Router struct {
	Rules         []Rule
	Fallback      interfaces.ModelProvider
	FallbackModel string
}

// Creates a Router with the given rules and fallback
func New(fallback interfaces.ModelProvider, fallbackModel string) *Router {
	return &Router{
		Fallback:      fallback,
		FallbackModel: fallbackModel,
	}
}

// Appends a routing rule
func (r *Router) AddRule(rule Rule) {
	r.Rules = append(r.Rules, rule)
}

func (r *Router) Route(prompt string) (interfaces.ModelProvider, string) {
	for _, rule := range r.Rules {
		if rule.Condition(prompt) {
			return rule.Provider, rule.Model
		}
	}
	return r.Fallback, r.FallbackModel
}

func IsShortPrompt(prompt string) bool {
	return len(strings.TrimSpace(prompt)) < 50
}

func IsLongPrompt(prompt string) bool {
	return len(strings.TrimSpace(prompt)) >= 50
}

func ContainsCode(prompt string) bool {
	lower := strings.ToLower(prompt)
	codeKeywords := []string{
		"code", "function", "debug", "error", "implement",
		"write a", "fix", "refactor", "algorithm",
	}
	for _, keyword := range codeKeywords {
		if strings.Contains(lower, keyword) {
			return true
		}
	}
	return false
}
