package handlers

import (
	"encoding/json"
	e "errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mohdjishin/SplitWise/helper/validate"
	"github.com/mohdjishin/SplitWise/internal/db"
	"github.com/mohdjishin/SplitWise/internal/errors"
	"github.com/mohdjishin/SplitWise/internal/middleware"
	"github.com/mohdjishin/SplitWise/internal/models"
	dto "github.com/mohdjishin/SplitWise/internal/models/dto"
	"github.com/mohdjishin/SplitWise/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CreateGroupWithBill(w http.ResponseWriter, r *http.Request) {

	var input dto.CreateGroupWithBillRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.LoggerInstance.Error("Error decoding request body", zap.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errors.ErrBadRequest)
		return
	}
	if err := validate.ValidateStruct(input); err != nil {
		logger.LoggerInstance.Error("Error validating request body", zap.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errors.ErrValidationFailed(err.Error()))
		return
	}

	userId := middleware.GetCurrentUserId(r)

	group := models.Group{
		Name:      input.GroupName,
		CreatedBy: uint(userId),
	}
	if err := db.GetDb().Create(&group).Error; err != nil {
		logger.LoggerInstance.Error("Failed to create group", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)

		return
	}

	groupMember := models.GroupMember{
		GroupID: group.ID,
		UserID:  uint(userId),
	}
	if err := db.GetDb().Create(&groupMember).Error; err != nil {
		logger.LoggerInstance.Error("Failed to add creator to group", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		return
	}

	bill := models.Bill{
		Name:    input.Bill.Name,
		Amount:  input.Bill.Amount,
		GroupID: group.ID,
	}
	if err := db.GetDb().Create(&bill).Error; err != nil {
		logger.LoggerInstance.Error("Failed to create bill", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		return
	}

	group.BillID = bill.ID
	if err := db.GetDb().Save(&group).Error; err != nil {
		logger.LoggerInstance.Error("Failed to update group with bill", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"groupId": group.ID,
		"billId":  bill.ID,
		"message": "Group and bill created successfully",
	})
}

func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	var group models.Group

	if err := db.GetDb().Where("id = ?", groupID).First(&group).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger.LoggerInstance.Error("Group not found", zap.Error(err))
		_ = json.NewEncoder(w).Encode(errors.ErrGroupNotFound)
		return
	}

	userId := middleware.GetCurrentUserId(r)
	if group.CreatedBy != uint(userId) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errors.ErrGroupNotFound)
		return
	}

	if err := db.GetDb().Delete(&group).Error; err != nil {
		logger.LoggerInstance.Error("Failed to delete group", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Group deleted"})
}

func ListOwnedGroups(w http.ResponseWriter, r *http.Request) {
	userId := middleware.GetCurrentUserId(r)
	var groups []models.Group

	err := db.GetDb().Table("groups").
		Select("groups.*, group_members.user_id").
		Joins("LEFT JOIN group_members ON group_members.group_id = groups.id").
		Where("group_members.user_id = ?", userId).
		Find(&groups).Error

	if err != nil {
		logger.LoggerInstance.Error("Failed to fetch groups", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		return
	}

	var groupList []struct {
		Group   models.Group         `json:"group"`
		Members []models.GroupMember `json:"members"`
	}

	for _, group := range groups {
		var members []models.GroupMember
		err := db.GetDb().Where("group_id = ?", group.ID).Find(&members).Error
		if err == nil {
			groupList = append(groupList, struct {
				Group   models.Group         `json:"group"`
				Members []models.GroupMember `json:"members"`
			}{
				Group:   group,
				Members: members,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(groupList)
}

func AddUsersToGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")

	var input dto.AddUsersToGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.LoggerInstance.Error("Invalid request payload", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errors.ErrBadRequest)
		return
	}
	if err := validate.ValidateStruct(input); err != nil {
		logger.LoggerInstance.Error("Error validating request body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errors.ErrValidationFailed(err.Error()))
		return
	}
	userId := middleware.GetCurrentUserId(r)
	var group models.Group
	if err := db.GetDb().Where("id = ? AND created_by = ?", groupID, userId).First(&group).Error; err != nil {
		logger.LoggerInstance.Error("Group not found", zap.Error(err))
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errors.ErrGroupNotFound)
		return
	}

	userList := map[string]uint{}
	missingUsers := []string{}
	for _, email := range input.UserEmailIds {
		var user models.User
		if err := db.GetDb().Where("email = ?", email).First(&user).Error; err != nil {
			missingUsers = append(missingUsers, email)
			continue
		}
		userList[email] = user.ID
	}

	if len(missingUsers) > 0 {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errors.ErrUsersNotFound(missingUsers))
		return
	}

	for _, userID := range input.UserEmailIds {
		var groupMember models.GroupMember
		if err := db.GetDb().Where("group_id = ? AND user_id = ?", groupID, userList[userID]).First(&groupMember).Error; err == nil { // TODO: change this. query once and check if the user is already a member
			continue
		}

		newGroupMember := models.GroupMember{
			GroupID: group.ID,
			UserID:  userList[userID],
		}
		if err := db.GetDb().Create(&newGroupMember).Error; err != nil {
			logger.LoggerInstance.Error("Failed to add user to group", zap.Error(err))
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
			return
		}
	}

	var groupMembers []models.GroupMember
	if err := db.GetDb().Where("group_id = ?", groupID).Find(&groupMembers).Error; err != nil {
		logger.LoggerInstance.Error("Failed to fetch group members", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrWhileFetchingMembers)
		return
	}

	var bill models.Bill
	if err := db.GetDb().Where("id = ?", group.BillID).First(&bill).Error; err != nil {
		logger.LoggerInstance.Error("Failed to fetch bill", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrWhileFetchingBill)
		return
	}

	numMembers := len(groupMembers)
	if numMembers > 0 {
		perUserSplitAmount := bill.Amount / float64(numMembers)
		group.TotalAmount = bill.Amount
		group.PerUserSplitAmount = perUserSplitAmount

		if err := db.GetDb().Save(&group).Error; err != nil {
			logger.LoggerInstance.Error("Failed to update group with total and per-user split amounts", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(errors.ErrInternalErrorWithMessage("Failed to update total and per-user split amounts"))
			return
		}

		for _, member := range groupMembers {
			member.SplitAmount = perUserSplitAmount
			if err := db.GetDb().Save(&member).Error; err != nil {
				logger.LoggerInstance.Error("Failed to update member split amount", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(errors.ErrInternalErrorWithMessage("Failed to update split amount for members"))
				return
			}
		}
	} else {
		group.TotalAmount = 0
		group.PerUserSplitAmount = 0
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Users added to group successfully"})
}

// ListMemberGroups
// @Summary List groups the user belongs to
// @Description Retrieves all groups associated with the authenticated user. Optionally filters the results by group status. If no status is provided, all groups will be returned.
// @Tags groups
// @Accept json
// @Produce json
// @Param status query string false "The status of the groups to filter by. Valid values are 'PENDING' or 'DONE'"
// @Success 200 {array} dto.ListMemberGroupsResponse "Successful response with the list of groups."
// @Failure 400 {object} errors.Error "Invalid status parameter provided."
// @Failure 500 {object} errors.Error "Internal server error."
// @Router /v1/groups/member-groups [get]
func ListMemberGroups(w http.ResponseWriter, r *http.Request) {
	userId := middleware.GetCurrentUserId(r)
	var groups []models.Group

	status := r.URL.Query().Get("status")

	if status != "PENDING" && status != "DONE" && status != "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errors.ErrInvalidQueryParameter("Invalid status for query parameter (status)"))
		return
	}
	query := db.GetDb().
		Joins("JOIN group_members ON groups.id = group_members.group_id").
		Where("group_members.user_id = ?", userId)

	if status != "" {
		query = query.Where("groups.status = ?", status)
	}

	err := query.Find(&groups).Error
	if err != nil {
		logger.LoggerInstance.Error("Failed to fetch groups", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		return
	}

	var groupList []dto.ListMemberGroupsResponse
	if len(groups) > 0 {
		var members []models.GroupMember
		err = db.GetDb().Where("group_id IN (?)", getGroupIDs(groups)).Find(&members).Error
		if err != nil {
			logger.LoggerInstance.Error("Failed to fetch group members", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
			return
		}

		groupMemberMap := make(map[uint][]models.GroupMember)
		for _, member := range members {
			groupMemberMap[member.GroupID] = append(groupMemberMap[member.GroupID], member)
		}

		for _, group := range groups {
			groupList = append(groupList, dto.ListMemberGroupsResponse{
				Group:   group,
				Members: groupMemberMap[group.ID],
			})
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(groupList)
}

func getGroupIDs(groups []models.Group) []uint {
	groupIDs := make([]uint, len(groups))
	for i, group := range groups {
		groupIDs[i] = group.ID
	}
	return groupIDs
}

// GetPendingPayments retrieves all pending payments associated with the authenticated user.
// @Summary Retrieve Pending Payments
// @Description Fetches all pending payments for the current user that have not been paid yet, including group ID, group name, bill ID, and amount owed.
// @Tags groups
// @Accept json
// @Produce json
// @Success 200 {object} dto.PendingPaymentsWithTotalResponse "Successful response containing the list of pending payments and total amount."
// @Failure 404 {object} errors.Error "No pending payments found for the user."
// @Failure 500 {object} errors.Error "Internal server error occurred while fetching pending payments."
// @Router /v1/groups/pending-payments [get]
func GetPendingPayments(w http.ResponseWriter, r *http.Request) {
	userId := middleware.GetCurrentUserId(r)
	var groupMembers []models.GroupMember

	if err := db.GetDb().Where("user_id = ? AND has_paid = ?", userId, false).Find(&groupMembers).Error; err != nil {
		if e.Is(err, gorm.ErrRecordNotFound) {
			logger.LoggerInstance.Warn("No pending payments found", zap.Float64("userId", userId))
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(errors.ErrNoPendingPayments)
		} else {
			logger.LoggerInstance.Error("Failed to fetch pending payments", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		}
		return
	}

	var pendingPayments []dto.PendingPayments
	totalAmount := 0.0

	for _, member := range groupMembers {
		var group models.Group
		if err := db.GetDb().Where("id = ?", member.GroupID).First(&group).Error; err != nil {
			logger.LoggerInstance.Error("Failed to fetch group", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
			return
		}

		var bill models.Bill
		if err := db.GetDb().Where("id = ?", group.BillID).First(&bill).Error; err != nil {
			logger.LoggerInstance.Error("Failed to fetch bill", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
			return
		}

		pendingPayments = append(pendingPayments, dto.PendingPayments{
			GroupID:   group.ID,
			GroupName: group.Name,
			BillID:    bill.ID,
			Amount:    bill.Amount,
		})
		totalAmount += bill.Amount
	}

	response := dto.PendingPaymentsWithTotalResponse{
		PendingPayments: pendingPayments,
		TotalAmount:     totalAmount,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
