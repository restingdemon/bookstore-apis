package models

import(
	"github.com/jinzhu/gorm"
	"github.com/restingdemon/go-mysql-mux/pkg/config"
	"github.com/google/uuid"
)

var Db *gorm.DB

type Book struct{
	gorm.Model
	Name string `json:"name"`
	Author string `json:"author"`
	Publication string `json:"publication"`
}

type User struct{
	gorm.Model
	First_name			string     `json:"first_name" validate:"required,min=2,max=100"`
	Last_name			string     `json:"last_name" validate:"required,min=2,max=100"`
	Password			string     `json:"password" validate:"required,min=6,max=100"`
	Email				string     `json:"email" validate:"email,required"`
	Phone				string     `json:"phone" validate:"required"`
	Token				string     `gorm:"size:2000" json:"token"`
	User_type			string     `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	Refresh_token		string     `json:"refresh_token"`
	User_id				string     `json:"user_id"`
}

func init(){
	config.Connect()
	Db = config.GetDB()
	Db.AutoMigrate(&Book{})
	Db.AutoMigrate(&User{})
}

func (b *Book) CreateBook() *Book{
	Db.NewRecord(b)
	Db.Create(&b)
	return b
}

func GetAllBooks() []Book{
	var Books []Book
	Db.Find(&Books)
	return Books
}

func GetBookById(Id int64) (*Book,*gorm.DB){
	var getBook Book
	Db:=Db.Where("ID=?",Id).Find(&getBook)
	return &getBook,Db
}

func DeleteBook(ID int64) Book{
	var book Book
	Db.Where("ID=?",ID).Delete(book)
	return book
}

func GetUserById(Id string) (*User,*gorm.DB){
	var getUser User
	Db:=Db.Where("User_id=?",Id).Find(&getUser)
	return &getUser,Db
}

func (b *User) CreateUser() *User{
	Db.NewRecord(b)
	Db.Create(&b)
	return b
}
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.User_id = uuid.New().String()
	return nil
}