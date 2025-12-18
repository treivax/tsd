// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

const (
	// MinConsumptions nombre minimum de consommations autorisées
	MinConsumptions = 1
)

// OnceConsumptionPolicy permet une seule consommation au total.
type OnceConsumptionPolicy struct{}

// NewOnceConsumptionPolicy crée une nouvelle politique de consommation unique.
func NewOnceConsumptionPolicy() *OnceConsumptionPolicy {
	return &OnceConsumptionPolicy{}
}

// CanConsume vérifie si le xuple peut être consommé (aucune consommation encore).
func (p *OnceConsumptionPolicy) CanConsume(xuple *Xuple, agentID string) bool {
	return xuple.Metadata.ConsumptionCount == 0
}

// OnConsumed marque le xuple comme complètement consommé.
func (p *OnceConsumptionPolicy) OnConsumed(xuple *Xuple, agentID string) bool {
	return true // Marquer comme complètement consommé
}

// Name retourne le nom de la politique.
func (p *OnceConsumptionPolicy) Name() string {
	return "once"
}

// PerAgentConsumptionPolicy permet une consommation par agent.
type PerAgentConsumptionPolicy struct{}

// NewPerAgentConsumptionPolicy crée une nouvelle politique par agent.
func NewPerAgentConsumptionPolicy() *PerAgentConsumptionPolicy {
	return &PerAgentConsumptionPolicy{}
}

// CanConsume vérifie si cet agent n'a pas encore consommé le xuple.
func (p *PerAgentConsumptionPolicy) CanConsume(xuple *Xuple, agentID string) bool {
	if xuple.Metadata.ConsumedBy == nil {
		return true
	}
	_, alreadyConsumed := xuple.Metadata.ConsumedBy[agentID]
	return !alreadyConsumed
}

// OnConsumed ne marque jamais comme complètement consommé.
func (p *PerAgentConsumptionPolicy) OnConsumed(xuple *Xuple, agentID string) bool {
	return false // Autres agents peuvent consommer
}

// Name retourne le nom de la politique.
func (p *PerAgentConsumptionPolicy) Name() string {
	return "per-agent"
}

// LimitedConsumptionPolicy permet un nombre limité de consommations.
type LimitedConsumptionPolicy struct {
	MaxConsumptions int
}

// NewLimitedConsumptionPolicy crée une nouvelle politique avec limite.
func NewLimitedConsumptionPolicy(maxConsumptions int) *LimitedConsumptionPolicy {
	if maxConsumptions <= 0 {
		maxConsumptions = MinConsumptions
	}
	return &LimitedConsumptionPolicy{
		MaxConsumptions: maxConsumptions,
	}
}

// CanConsume vérifie si la limite n'est pas atteinte.
func (p *LimitedConsumptionPolicy) CanConsume(xuple *Xuple, agentID string) bool {
	return xuple.Metadata.ConsumptionCount < p.MaxConsumptions
}

// OnConsumed marque comme complètement consommé si limite atteinte.
func (p *LimitedConsumptionPolicy) OnConsumed(xuple *Xuple, agentID string) bool {
	return xuple.Metadata.ConsumptionCount >= p.MaxConsumptions
}

// Name retourne le nom de la politique.
func (p *LimitedConsumptionPolicy) Name() string {
	return "limited"
}
