package router

import (
	"inference-gateway/interfaces"
	"strings"
)

// Defines a single routing decision
type Rule struct {
	Name         string
	Condition    func(prompt string) bool
	Provider     interfaces.ModelProvider
	Model        string
	ProviderName string
}

// Holds all routing rules and fallback provider
type Router struct {
	Rules                []Rule
	Fallback             interfaces.ModelProvider
	FallbackModel        string
	FallbackProviderName string
}

// Creates a Router with the given rules and fallback
func New(fallback interfaces.ModelProvider, fallbackModel string, fallbackProviderName string) *Router {
	return &Router{
		Fallback:             fallback,
		FallbackModel:        fallbackModel,
		FallbackProviderName: fallbackProviderName,
	}
}

// Appends a routing rule
func (r *Router) AddRule(rule Rule) {
	r.Rules = append(r.Rules, rule)
}

func (r *Router) Route(prompt string) (interfaces.ModelProvider, string, string) {
	for _, rule := range r.Rules {
		if rule.Condition(prompt) {
			return rule.Provider, rule.Model, rule.ProviderName
		}
	}
	return r.Fallback, r.FallbackModel, r.FallbackProviderName
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
