<<<<<<< HEAD
package message

type User struct {
	Uid          int
	Uname        string
	Admin        bool
	Urank        int
	MobileVerify bool
	Medal        *Medal
	GuardLevel   int
}

type Medal struct {
	Name     string
	Level    int
	Color    int
	UpRoomId int
	UpUid    int
	UpName   string
}
=======
package message

type User struct {
	Uid          int
	Uname        string
	Admin        bool
	Urank        int
	MobileVerify bool
	Medal        *Medal
	GuardLevel   int
}

type Medal struct {
	Name     string
	Level    int
	Color    int
	UpRoomId int
	UpUid    int
	UpName   string
}
>>>>>>> e20f45e8c9dc9dc6e202d459cd56a928afdd5f95
