package repositories

import (
	"fmt"
	"gokg/gomvc/models"
	"sync"
	"time"
)

var Audits = []models.AuditLog{}

type AuditInterface interface {
	AuditsAll() []models.AuditLog
	CreateInsertEvent(userID int, newState models.IAuditable)
	CreateUpdateEvent(userID int, oldState models.IAuditable, newState models.IAuditable)
	CreateDeleteEvent(userID int, oldState models.IAuditable)
}

type AuditBase struct {
	auditInterface AuditInterface
}

func (auditBase *AuditBase) Init(auditClass AuditInterface) {
	auditBase.auditInterface = auditClass
}

type AuditAbstract struct {
	_itemCount int
	mu         sync.Mutex
}

func NewAuditAbstract() *AuditAbstract {
	roleAbstract := &AuditAbstract{}
	roleAbstract._itemCount = len(Audits)
	return roleAbstract
}

func (acl *AuditAbstract) AuditsAll() []models.AuditLog {
	fmt.Println("all audits abstract requested!")
	return Audits
}

func (acl *AuditAbstract) InsertAudit(audit *models.AuditLog) {
	// Lock the mutex before accessing _userCounter
	acl.mu.Lock()
	defer acl.mu.Unlock()
	audit.ID = acl._itemCount + 1
	acl._itemCount++
	Audits = append(Audits, *audit)
	fmt.Printf("abstract audit appended: %v", audit)
	fmt.Println("")
}

func (auditor *AuditAbstract) InsertAuditWithDetails(userID int,
	oldState models.IAuditable, newState models.IAuditable, action string) {
	var oldStateStr string
	if oldState != nil {
		oldStateStr = oldState.GetObjectStr()
	} else {
		oldStateStr = ""
	}

	var newStateStr string
	if newState != nil {
		newStateStr = newState.GetObjectStr()
	} else {
		newStateStr = ""
	}

	var IDStr string
	if oldState != nil {
		IDStr = oldState.GetID()
	} else {
		IDStr = ""
	}
	if IDStr == "" {
		if newState != nil {
			IDStr = newState.GetID()
		} else {
			IDStr = ""
		}
	}

	var objectName string
	if newState != nil {
		objectName = newState.GetObjectName()
	} else {
		objectName = ""
	}

	auditLog := &models.AuditLog{
		UserID:     userID,
		Action:     action,
		ObjectName: objectName,
		Timestamp:  time.Now(),
		Details:    fmt.Sprintf("%s with ID %v", action, IDStr),
		OldState:   oldStateStr,
		NewState:   newStateStr,
	}
	auditor.InsertAudit(auditLog)
}

func (acl *AuditAbstract) CreateInsertEvent(userID int, newState models.IAuditable) {
	acl.InsertAuditWithDetails(userID, newState, nil, "Inserted")
}
func (acl *AuditAbstract) CreateUpdateEvent(userID int, oldState models.IAuditable, newState models.IAuditable) {
	acl.InsertAuditWithDetails(userID, oldState, newState, "Updated")
}
func (acl *AuditAbstract) CreateDeleteEvent(userID int, oldState models.IAuditable) {
	acl.InsertAuditWithDetails(userID, oldState, oldState, "Deleted")
}

type AuditCustom struct {
	_itemCount int
	mu         sync.Mutex
}

func NewAuditCustom() *AuditCustom {
	auditCustom := &AuditCustom{}
	auditCustom._itemCount = len(Audits)
	return auditCustom
}

func (acl *AuditCustom) AuditsAll() []models.AuditLog {
	fmt.Println("all audits custom requested!")
	return Audits
}

func (acl *AuditCustom) InsertAudit(audit *models.AuditLog) {
	// Lock the mutex before accessing _userCounter
	acl.mu.Lock()
	defer acl.mu.Unlock()
	audit.ID = acl._itemCount + 1
	acl._itemCount++
	Audits = append(Audits, *audit)
	fmt.Printf("custom audit appended: %v", audit)
	fmt.Println("")
}

func (auditor *AuditCustom) InsertAuditWithDetails(userID int, oldState models.IAuditable, newState models.IAuditable, action string) {
	auditLog := &models.AuditLog{
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now(),
		Details:   fmt.Sprintf("%s with ID %v", action, oldState.GetID()),
		OldState:  oldState.GetObjectStr(),
		NewState:  newState.GetObjectStr(),
	}
	auditor.InsertAudit(auditLog)
}
