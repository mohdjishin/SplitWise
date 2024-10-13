package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mohdjishin/SplitWise/helper/validate"
	"github.com/mohdjishin/SplitWise/internal/db"
	"github.com/mohdjishin/SplitWise/internal/errors"
	"github.com/mohdjishin/SplitWise/internal/middleware"
	"github.com/mohdjishin/SplitWise/internal/models"
	dto "github.com/mohdjishin/SplitWise/internal/models/dto"
	log "github.com/mohdjishin/SplitWise/logger"
	"go.uber.org/zap"
)

// CreateGroupWithBill handles creating a group with an associated bill
// @Summary Create a new group with an associated bill
// @Description Creates a group with the specified name and an associated bill, then adds the user as a member of the group.
// @Tags groups
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body dto.CreateGroupWithBillRequest true "CreateGroupWithBillRequest details"
// @Success 201 {object} dto.CreateGroupWithBillResponse
// @Failure 400 {object} errors.Error "Bad Request"
// @Failure 500 {object} errors.Error "Internal Server Error"
// @Router /v1/groups/ [post]
func CreateGroupWithBill(w http.ResponseWriter, r *http.Request) {

	var input dto.CreateGroupWithBillRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Error("Error decoding request body", zap.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errors.ErrBadRequest)
		return
	}
	log.Debug("CreateGroupWithBill request", zap.Any("request", input))
	if err := validate.ValidateStruct(input); err != nil {
		log.Error("Error validating request body", zap.Any("error", err))
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
		log.Error("Failed to create group", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)

		return
	}

	groupMember := models.GroupMember{
		GroupID: group.ID,
		UserID:  uint(userId),
	}
	if err := db.GetDb().Create(&groupMember).Error; err != nil {
		log.Error("Failed to add creator to group", zap.Error(err))
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
		log.Error("Failed to create bill", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		return
	}

	group.BillID = bill.ID
	if err := db.GetDb().Save(&group).Error; err != nil {
		log.Error("Failed to update group with bill", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(dto.CreateGroupWithBillResponse{
		GroupID: group.ID,
		BillID:  bill.ID,
		Message: "Group and bill created successfully",
	})
}

// DeleteGroup handles deleting a specified group
// @Summary Delete a group by ID (NOT NEEDED AS OF NOW)
// @Description Deletes a group identified by the specified group ID if the user is the creator of the group. (owner only can do this operation)
// @Tags groups
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID of the group to be deleted"
// @Success 200 {object} dto.DeleteGroupResponse
// @Failure 404 {object} errors.Error "Group Not Found"
// @Failure 500 {object} errors.Error "Internal Server Error"
// @Router /v1/groups/{id} [delete]
func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	var group models.Group
	log.Debug("DeleteGroup request", zap.Any("groupID", groupID))
	if err := db.GetDb().Where("id = ?", groupID).First(&group).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Error("Group not found", zap.Error(err))
		_ = json.NewEncoder(w).Encode(errors.ErrGroupNotFound)
		return
	}

	userId := middleware.GetCurrentUserId(r)
	if group.CreatedBy != uint(userId) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errors.ErrGroupNotFound)
		return
	}
	log.Debug("DeleteGroup request", zap.Any("groupID and userId", []any{groupID, userId}))
	if err := db.GetDb().Delete(&group).Error; err != nil {
		log.Error("Failed to delete group", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dto.DeleteGroupResponse{Message: "Group deleted"})
}

// ListOwnedGroups handles fetching groups owned by the current user
// @Summary List groups owned by the user
// @Description Fetches and returns a list of groups that are owned by the current user, including group members.
// @Tags groups
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} dto.ListOwnedGroupsResponse "List of groups owned by the user"
// @Failure 500 {object} errors.Error "Internal Server Error"
// @Router /v1/groups/owned [get]
func ListOwnedGroups(w http.ResponseWriter, r *http.Request) {
	userId := middleware.GetCurrentUserId(r)
	var groups []models.Group
	log.Debug("ListOwnedGroups request", zap.Any("userId", userId))
	err := db.GetDb().Table("groups").
		Select("groups.*, group_members.user_id").
		Joins("LEFT JOIN group_members ON group_members.group_id = groups.id").
		Where("group_members.user_id = ?", userId).
		Find(&groups).Error

	if err != nil {
		log.Error("Failed to fetch groups", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		return
	}

	var groupList []dto.ListOwnedGroupsResponse
	for _, group := range groups {
		var members []models.GroupMember
		err := db.GetDb().Where("group_id = ?", group.ID).Find(&members).Error
		if err == nil {
			groupList = append(groupList, dto.ListOwnedGroupsResponse{
				Group:   group,
				Members: members,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(groupList)
}

// AddUsersToGroup handles adding users to a specified group
// @Summary Add members to a group
// @Description Adds members identified by their email addresses to a group if the user is the creator of the group.
// @Tags groups
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID of the group to which members will be added"
// @Param request body dto.AddUsersToGroupRequest true "List of user email IDs to add to the group"
// @Success 200 {object} dto.AddUsersToGroupResponse "success message"
// @Failure 400 {object} errors.Error "Bad Request"
// @Failure 404 {object} errors.Error "Group Not Found or Users Not Found"
// @Failure 500 {object} errors.Error "Internal Server Error"
// @Router /v1/groups/{id}/users [post]
func AddUsersToGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	log.Debug("AddUsersToGroup request", zap.Any("groupID", groupID))
	var input dto.AddUsersToGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Error("Invalid request payload", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errors.ErrBadRequest)
		return
	}
	log.Debug("AddUsersToGroup request", zap.Any("request", input))

	if err := validate.ValidateStruct(input); err != nil {
		log.Error("Error validating request body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errors.ErrValidationFailed(err.Error()))
		return
	}
	userId := middleware.GetCurrentUserId(r)
	log.Debug("AddUsersToGroup request", zap.Any("userId", userId))
	var group models.Group
	if err := db.GetDb().Where("id = ? AND created_by = ?", groupID, userId).First(&group).Error; err != nil {
		log.Error("Group not found", zap.Error(err))
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errors.ErrGroupNotFound)
		return
	}

	userList := map[string]uint{}
	missingUsers := []string{}
	existingUsers := []string{}
	for _, email := range input.UserEmailIds {
		var user models.User
		if err := db.GetDb().Where("email = ?", email).First(&user).Error; err != nil {
			missingUsers = append(missingUsers, email)
			continue
		}
		userList[email] = user.ID
		var groupMember models.GroupMember
		if err := db.GetDb().Where("group_id = ? AND user_id = ?", groupID, user.ID).First(&groupMember).Error; err == nil {
			existingUsers = append(existingUsers, email)
		}
	}

	if len(missingUsers) > 0 {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errors.ErrUsersNotFound(missingUsers))
		return
	}
	if len(existingUsers) > 0 {
		w.WriteHeader(http.StatusConflict) // HTTP 409 Conflict
		_ = json.NewEncoder(w).Encode(errors.ErrUsersAlreadyExists(existingUsers))
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
			log.Error("Failed to add user to group", zap.Error(err))
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
			return
		}
	}

	var groupMembers []models.GroupMember
	if err := db.GetDb().Where("group_id = ?", groupID).Find(&groupMembers).Error; err != nil {
		log.Error("Failed to fetch group members", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrWhileFetchingMembers)
		return
	}

	var bill models.Bill
	if err := db.GetDb().Where("id = ?", group.BillID).First(&bill).Error; err != nil {
		log.Error("Failed to fetch bill", zap.Error(err))
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
			log.Error("Failed to update group with total and per-user split amounts", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(errors.ErrInternalErrorWithMessage("Failed to update total and per-user split amounts"))
			return
		}

		for _, member := range groupMembers {
			member.SplitAmount = perUserSplitAmount
			if err := db.GetDb().Save(&member).Error; err != nil {
				log.Error("Failed to update member split amount", zap.Error(err))
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
	_ = json.NewEncoder(w).Encode(dto.AddUsersToGroupResponse{Message: "Users added to group successfully"})
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
	log.Debug("ListMemberGroups request", zap.Any("userId", userId))

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
		log.Error("Failed to fetch groups", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		return
	}
	log.Debug("ListMemberGroups request", zap.Any("groups", groups))
	log.Debug("len(groups)", zap.Any("len(groups)", len(groups)))

	var groupList []dto.ListMemberGroupsResponse
	if len(groups) > 0 {
		var members []models.GroupMember
		err = db.GetDb().Where("group_id IN (?)", getGroupIDs(groups)).Find(&members).Error
		if err != nil {
			log.Error("Failed to fetch group members", zap.Error(err))
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
