package service

import (
	"errors"
	"fmt"
	"gin-boilerplate/data/request"
	"gin-boilerplate/data/response"
	"gin-boilerplate/exception"
	"gin-boilerplate/helper"
	"gin-boilerplate/model"
	"gin-boilerplate/repository"
	"gin-boilerplate/utils"
	"mime/multipart"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/tealeg/xlsx"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	Validate       *validator.Validate
}

func NewUserServiceImpl(userRepository repository.UserRepository, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		Validate:       validate,
	}
}

func (service *UserServiceImpl) FindAll(filters map[string]string) []response.UsersResponse {
	result := service.UserRepository.FindAll(filters)

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

func (service *UserServiceImpl) Create(user request.CreateUsersRequest) {
	err := service.Validate.Struct(user)
	helper.ErrorPanic(err)

	hashedPassword, err := utils.HashPassword(user.Password)
	helper.ErrorPanic(err)

	dataset := model.Users{
		Username: user.Username,
		Email:    user.Email,
		Password: hashedPassword,
	}
	service.UserRepository.Save(dataset)
}

func (service *UserServiceImpl) FindById(userId int) response.UsersResponse {
	dataset, err := service.UserRepository.FindById(userId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	Response := response.UsersResponse{
		Id:       dataset.Id,
		Username: dataset.Username,
		Email:    dataset.Email,
	}
	return Response
}

func (service *UserServiceImpl) Update(user request.UpdateUsersRequest) {
	err := service.Validate.Struct(user)
	helper.ErrorPanic(err)

	dataset, err := service.UserRepository.FindById(user.Id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	if user.Password != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		helper.ErrorPanic(err)
		dataset.Password = hashedPassword
	}
	dataset.Username = user.Username
	dataset.Email = user.Email
	service.UserRepository.Update(dataset)
}

func (service *UserServiceImpl) BulkDelete(userIds []int) error {
	if len(userIds) == 0 {
		return errors.New("user not found")
	}
	service.UserRepository.BulkDelete(userIds)
	
	return nil
}

func (service *UserServiceImpl) Export() (string, error) {
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

	users := service.UserRepository.FindAll(nil) // Menggunakan nil karena tidak ada filter pada contoh ini

	for _, user := range users {
		dataRow := sheet.AddRow()
		dataRow.AddCell().SetInt(int(user.Id))
		dataRow.AddCell().Value = user.Username
		dataRow.AddCell().Value = user.Email
		dataRow.AddCell().Value = user.CreatedAt.Format("2006-01-02")
	}

	// Save the Excel file
	timestamp := time.Now().Format("2006-01-02_150405")
	filePath := fmt.Sprintf("user_%s.xlsx", timestamp)

	err = file.Save(filePath)
	if err != nil {
		return "", err
	}
	// return nil
	return filePath, nil
}

func (service *UserServiceImpl) Import(file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	xlFile, err := xlsx.OpenReaderAt(src, file.Size)
	if err != nil {
		return err
	}

	sheet := xlFile.Sheets[0]

	// Create channels for error handling and synchronization
	errorChan := make(chan error)
	wg := sync.WaitGroup{}

	for rowIndex, row := range sheet.Rows {
		if rowIndex == 0 {
			continue
		}

		wg.Add(1)

		go func(rowIndex int, row *xlsx.Row) {
			defer wg.Done()

			usernameCell := row.Cells[0]
			emailCell := row.Cells[1]

			if usernameCell.String() == "" || emailCell.String() == "" {
				errorChan <- fmt.Errorf("empty username or email cell found in row %d", rowIndex)
				return
			}

			existingUser, err := service.UserRepository.FindByEmail(emailCell.String())
			if err != nil {
				errorChan <- fmt.Errorf("error checking existing email in row %d: %v", rowIndex, err)
				return
			}
			if existingUser.Id != 0 {
				errorChan <- fmt.Errorf("email '%s' already taken in row %d", emailCell.String(), rowIndex)
				return
			}

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
			service.UserRepository.Save(user)

		}(rowIndex, row)
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	// Handle errors from goroutines
	for err := range errorChan {
		return err
	}

	return nil
}
