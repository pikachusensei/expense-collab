package repository

import "github.com/shreyansh/expense-go-collab-backend/internal/model"

type UserRepository interface {
	CreateUser(user *model.User) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUserByID(id int) (*model.User, error)
	GetAllUsers() ([]*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	DeleteUser(id int) error
}

type GroupRepository interface {
	CreateGroup(group *model.Group) (*model.Group, error)
	GetGroupByID(id int) (*model.Group, error)
	GetAllGroups() ([]*model.Group, error)
	GetGroupsByUserID(userID int) ([]*model.Group, error)
	UpdateGroup(group *model.Group) (*model.Group, error)
	DeleteGroup(id int) error
}

type GroupMemberRepository interface {
	AddMember(member *model.GroupMember) (*model.GroupMember, error)
	RemoveMember(groupID, userID int) error
	GetGroupMembers(groupID int) ([]*model.GroupMember, error)
	GetUserGroups(userID int) ([]*model.GroupMember, error)
	IsMember(groupID, userID int) (bool, error)
}

type ExpenseRepository interface {
	CreateExpense(expense *model.Expense) (*model.Expense, error)
	GetExpenseByID(id int) (*model.Expense, error)
	GetExpensesByGroupID(groupID int) ([]*model.Expense, error)
	GetExpensesByUserID(userID int) ([]*model.Expense, error)
	UpdateExpense(expense *model.Expense) (*model.Expense, error)
	DeleteExpense(id int) error
}

type ExpenseSplitRepository interface {
	CreateSplit(split *model.ExpenseSplit) (*model.ExpenseSplit, error)
	GetSplitsByExpenseID(expenseID int) ([]*model.ExpenseSplit, error)
	GetSplitsByUserID(userID int) ([]*model.ExpenseSplit, error)
	DeleteSplitsByExpenseID(expenseID int) error
	UpdateSplit(split *model.ExpenseSplit) (*model.ExpenseSplit, error)
}

type BalanceRepository interface {
	GetBalance(userID, groupID int) (float64, error)
	GetUserBalanceInGroup(userID, groupID int) (float64, error)
	GetGroupBalances(groupID int) (map[int]float64, error)
	CalculateBalances(groupID int) error
}

type SettlementRepository interface {
	CreateSettlement(settlement *model.Settlement) (*model.Settlement, error)
	GetSettlementByID(id int) (*model.Settlement, error)
	GetSettlementsByGroupID(groupID int) ([]*model.Settlement, error)
	GetSettlementsByUserID(userID int) ([]*model.Settlement, error)
	GetAllSettlements() ([]*model.Settlement, error)
}
