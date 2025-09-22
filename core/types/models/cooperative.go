package models

import (
	"time"
)

// CooperativeCredits represents credits for a node in the cooperative system
type CooperativeCredits struct {
	NodeID             string    `json:"node_id" gorm:"primaryKey;type:uuid"`
	Balance            float64   `json:"balance" gorm:"not null;default:0"`
	LastTransactionAt  time.Time `json:"last_transaction_at"`
	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relationships
	Node *Node `json:"node,omitempty" gorm:"foreignKey:NodeID"`
}

// CooperativeTransaction represents a transaction in the cooperative system
type CooperativeTransaction struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	FromNodeID  string    `json:"from_node_id" gorm:"not null;index"`
	ToNodeID    string    `json:"to_node_id" gorm:"not null;index"`
	Amount      float64   `json:"amount" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"`
	Description string    `json:"description"`
	Status      string    `json:"status" gorm:"not null;default:'pending'"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relationships
	FromNode *Node `json:"from_node,omitempty" gorm:"foreignKey:FromNodeID"`
	ToNode   *Node `json:"to_node,omitempty" gorm:"foreignKey:ToNodeID"`
}

// CooperativeTransactionType represents the possible types of cooperative transactions
type CooperativeTransactionType string

const (
	TransactionTypeServiceReward    CooperativeTransactionType = "service_reward"
	TransactionTypeResourceUsage    CooperativeTransactionType = "resource_usage"
	TransactionTypeParticipation    CooperativeTransactionType = "participation"
	TransactionTypeTransfer         CooperativeTransactionType = "transfer"
	TransactionTypePenalty          CooperativeTransactionType = "penalty"
	TransactionTypeBonus            CooperativeTransactionType = "bonus"
)

// IsValid checks if the transaction type is valid
func (t CooperativeTransactionType) IsValid() bool {
	switch t {
	case TransactionTypeServiceReward, TransactionTypeResourceUsage,
		 TransactionTypeParticipation, TransactionTypeTransfer,
		 TransactionTypePenalty, TransactionTypeBonus:
		return true
	default:
		return false
	}
}

// CooperativeTransactionStatus represents the possible statuses of a cooperative transaction
type CooperativeTransactionStatus string

const (
	TransactionStatusPending   CooperativeTransactionStatus = "pending"
	TransactionStatusCompleted CooperativeTransactionStatus = "completed"
	TransactionStatusFailed    CooperativeTransactionStatus = "failed"
	TransactionStatusCancelled CooperativeTransactionStatus = "cancelled"
)

// IsValid checks if the transaction status is valid
func (s CooperativeTransactionStatus) IsValid() bool {
	switch s {
	case TransactionStatusPending, TransactionStatusCompleted,
		 TransactionStatusFailed, TransactionStatusCancelled:
		return true
	default:
		return false
	}
}

// GovernanceProposal represents a governance proposal in the cooperative system
type GovernanceProposal struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	ProposerID  string    `json:"proposer_id" gorm:"not null;index"`
	Status      string    `json:"status" gorm:"not null;default:'active'"`
	VotesFor    int       `json:"votes_for" gorm:"default:0"`
	VotesAgainst int      `json:"votes_against" gorm:"default:0"`
	VotesAbstain int      `json:"votes_abstain" gorm:"default:0"`
	StartDate   time.Time `json:"start_date" gorm:"not null"`
	EndDate     time.Time `json:"end_date" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relationships
	Proposer *Node `json:"proposer,omitempty" gorm:"foreignKey:ProposerID"`
}

// GovernanceProposalStatus represents the possible statuses of a governance proposal
type GovernanceProposalStatus string

const (
	ProposalStatusActive   GovernanceProposalStatus = "active"
	ProposalStatusVoting   GovernanceProposalStatus = "voting"
	ProposalStatusPassed   GovernanceProposalStatus = "passed"
	ProposalStatusRejected GovernanceProposalStatus = "rejected"
	ProposalStatusExpired  GovernanceProposalStatus = "expired"
)

// IsValid checks if the proposal status is valid
func (s GovernanceProposalStatus) IsValid() bool {
	switch s {
	case ProposalStatusActive, ProposalStatusVoting,
		 ProposalStatusPassed, ProposalStatusRejected, ProposalStatusExpired:
		return true
	default:
		return false
	}
}

// GovernanceVote represents a vote on a governance proposal
type GovernanceVote struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProposalID string    `json:"proposal_id" gorm:"not null;index"`
	VoterID    string    `json:"voter_id" gorm:"not null;index"`
	Vote       string    `json:"vote" gorm:"not null"`
	Reason     string    `json:"reason"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relationships
	Proposal *GovernanceProposal `json:"proposal,omitempty" gorm:"foreignKey:ProposalID"`
	Voter    *Node               `json:"voter,omitempty" gorm:"foreignKey:VoterID"`
}

// GovernanceVoteType represents the possible types of governance votes
type GovernanceVoteType string

const (
	VoteTypeYes     GovernanceVoteType = "yes"
	VoteTypeNo      GovernanceVoteType = "no"
	VoteTypeAbstain GovernanceVoteType = "abstain"
)

// IsValid checks if the vote type is valid
func (v GovernanceVoteType) IsValid() bool {
	switch v {
	case VoteTypeYes, VoteTypeNo, VoteTypeAbstain:
		return true
	default:
		return false
	}
}

// NodeReputation represents the reputation score for a node
type NodeReputation struct {
	NodeID      string    `json:"node_id" gorm:"primaryKey;type:uuid"`
	Score       float64   `json:"score" gorm:"not null;default:0"`
	TrustLevel  string    `json:"trust_level" gorm:"not null;default:'unknown'"`
	LastUpdated time.Time `json:"last_updated" gorm:"autoUpdateTime"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relationships
	Node *Node `json:"node,omitempty" gorm:"foreignKey:NodeID"`
}

// NodeTrustLevel represents the possible trust levels for a node
type NodeTrustLevel string

const (
	TrustLevelUnknown NodeTrustLevel = "unknown"
	TrustLevelLow     NodeTrustLevel = "low"
	TrustLevelMedium  NodeTrustLevel = "medium"
	TrustLevelHigh    NodeTrustLevel = "high"
	TrustLevelVeryHigh NodeTrustLevel = "very_high"
)

// IsValid checks if the trust level is valid
func (t NodeTrustLevel) IsValid() bool {
	switch t {
	case TrustLevelUnknown, TrustLevelLow, TrustLevelMedium,
		 TrustLevelHigh, TrustLevelVeryHigh:
		return true
	default:
		return false
	}
}

// GetTrustLevel returns the trust level based on the reputation score
func (r *NodeReputation) GetTrustLevel() NodeTrustLevel {
	switch {
	case r.Score >= 9.0:
		return TrustLevelVeryHigh
	case r.Score >= 7.0:
		return TrustLevelHigh
	case r.Score >= 5.0:
		return TrustLevelMedium
	case r.Score >= 3.0:
		return TrustLevelLow
	default:
		return TrustLevelUnknown
	}
}

// TableName returns the table name for the CooperativeCredits model
func (CooperativeCredits) TableName() string {
	return "cooperative_credits"
}

// TableName returns the table name for the CooperativeTransaction model
func (CooperativeTransaction) TableName() string {
	return "cooperative_transactions"
}

// TableName returns the table name for the GovernanceProposal model
func (GovernanceProposal) TableName() string {
	return "governance_proposals"
}

// TableName returns the table name for the GovernanceVote model
func (GovernanceVote) TableName() string {
	return "governance_votes"
}

// TableName returns the table name for the NodeReputation model
func (NodeReputation) TableName() string {
	return "node_reputation"
}
