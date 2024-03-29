package campaign

import (
	"bwastartup/api/user"
	"time"

	"github.com/leekchan/accounting"
)

type Campaign struct {
	ID               int
	UserID           int
	Name             string
	ShortDescription string
	Description      string
	Perks            string
	BackerCount      int
	GoalAmount       int
	CurrentAmount    int
	Slug             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CampaignImages   []CampaignImage
	User 						 user.User
}

type CampaignImage struct {
	ID         int
	CampaignID int
	FileName   string
	IsPrimary  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (c Campaign) GoalAmountFormatIDR() string {
	ac := accounting.Accounting{Symbol: "Rp.", Precision: 2, Thousand: ".", Decimal: ","}
	return ac.FormatMoney(c.GoalAmount)
}