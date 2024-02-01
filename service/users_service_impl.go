package service

import (
	"fmt"
	"gin-boilerplate/data/request"
	"gin-boilerplate/data/response"
	"gin-boilerplate/helper"
	"gin-boilerplate/model"
	"gin-boilerplate/repository"
	"gin-boilerplate/utils"
	"mime/multipart"
	"time"

	"github.com/tealeg/xlsx"
)

type UsersServiceImpl struct {
	UsersRepository repository.UsersRepository
}

func NewUsersServiceImpl(usersRepository repository.UsersRepository) UsersService {
	return &UsersServiceImpl{
		UsersRepository: usersRepository,
	}
}

func (u *UsersServiceImpl) FindAll(filters map[string]string) []response.UsersResponse {
	result := u.UsersRepository.FindAll(filters)

	var users []response.UsersResponse
	for _, value := range result {
		user := response.UsersResponse{
			Id:       value.Id,
			Username: value.Username,
			Email:    value.Email,
		}
		users = append(users, user)
	}

	return users
}

func (a *UsersServiceImpl) Create(users request.CreateUsersRequest) {
	
	hashedPassword, err := utils.HashPassword(users.Password)
	helper.ErrorPanic(err)
	
	newUser := model.Users{
		Username: users.Username,
		Email:    users.Email,
		Password: hashedPassword,
	}
	a.UsersRepository.Save(newUser)
}

func (u *UsersServiceImpl) FindById(userId int) response.UsersResponse {
	userData, err := u.UsersRepository.FindById(userId)
	helper.ErrorPanic(err)
	
	userResponse := response.UsersResponse{
		Id:   userData.Id,
		Username: userData.Username,
		Email: userData.Email,
	}
	return userResponse
}


func (u *UsersServiceImpl) Update(users request.UpdateUsersRequest) {
	userData, err := u.UsersRepository.FindById(users.Id)
	helper.ErrorPanic(err)
	if users.Password != "" {
		hashedPassword, err := utils.HashPassword(users.Password)
		helper.ErrorPanic(err)
		userData.Password = hashedPassword
	}
	userData.Username = users.Username
	userData.Email = users.Email
	u.UsersRepository.Update(userData)
}


// func (u *UsersServiceImpl) Delete(userId int) {
// 	u.UsersRepository.Delete(userId)
// }

func (u *UsersServiceImpl) BulkDelete(userIds []int) {
	if len(userIds) == 0 {
		return
	}
	u.UsersRepository.BulkDelete(userIds)
}

func (u *UsersServiceImpl) Export()  (string, error) {
	// Create a new Excel file
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Users")
	helper.ErrorPanic(err)

	headers := []string{"Id", "Username", "Email", "CreatedAt"} // Sesuaikan header dengan struktur model.Users
	headerRow := sheet.AddRow()
	for _, header := range headers {
		cell := headerRow.AddCell()
		cell.Value = header
	}

	users := u.UsersRepository.FindAll(nil) // Menggunakan nil karena tidak ada filter pada contoh ini

	for _, user := range users {
		dataRow := sheet.AddRow()
		dataRow.AddCell().SetInt(int(user.Id))
		dataRow.AddCell().Value = user.Username
		dataRow.AddCell().Value = user.Email
		dataRow.AddCell().Value = user.CreatedAt.Format("2006-01-02")
	}

	// Save the Excel file
	filePath := "user.xlsx"
	err = file.Save(filePath)
	if err != nil {
		return "", err
	}
	// return nil
	return filePath, nil
}

func (u *UsersServiceImpl) Import(file *multipart.FileHeader) error {

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Read the Excel file
	xlFile, err := xlsx.OpenReaderAt(src, file.Size)
	if err != nil {
		return err
	}

	// Assume that the first sheet contains user data
	sheet := xlFile.Sheets[0]

    for rowIndex, row := range sheet.Rows {
        // Skip the first row (header)
        if rowIndex == 0 {
            continue
        }

        // Check if username cell (index 0) or email cell (index 1) is empty
        usernameCell := row.Cells[0]
        emailCell := row.Cells[1]

        if usernameCell.String() == "" || emailCell.String() == "" {
            // Log or handle the empty username or email cell as needed
            fmt.Printf("Empty username or email cell found in row %d\n", rowIndex)
            continue
        }

        // Check if email already exists in the database
        existingUser, err := u.UsersRepository.FindByEmail(emailCell.String())
         if err != nil {
            return fmt.Errorf("Error checking existing email in row %d: %v", rowIndex, err)
        }
        if existingUser.Id != 0 {
            return fmt.Errorf("Email '%s' already taken in row %d", emailCell.String(), rowIndex)
        }

        // Assuming the data is in the format: ID, Username, Email, CreatedAt
        var user model.Users

        for i, cell := range row.Cells {
            switch i {
            case 0:
                user.Username = cell.String()
            case 1:
                user.Email = cell.String()
            case 2:
                dateStr := cell.String()
                user.CreatedAt, _ = time.Parse("2006-01-02", dateStr)
            }
        }
        // Save or process the user data as needed
        u.UsersRepository.Save(user)
    }

    return nil
}
