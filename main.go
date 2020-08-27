package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"./dao"
	"./game"
	"./session"
)

const (
	regularID          = "^[0-9]{1,12}$"
	regularUsername    = "^[a-zA-Z][a-zA-Z0-9_]{4,15}$"
	regularPassword    = "^[a-zA-Z0-9_]{4,25}$"
	regularPasswordEnc = "^[a-zA-Z0-9_]{30,40}$"
	regularMkey        = "^[a-zA-Z0-9_]{15,20}$"
	regularPkey        = "^[a-zA-Z0-9_,]{15,64}$"
	regularPhone       = "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	regularName        = "^[a-zA-Z0-9_\u4e00-\u9fa5]{2,8}$"
	regularIP          = "^((25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9]).){3}(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])$"
	// regularName        = "^[a-zA-Z0-9_]{2,15}$"
	// regularName        = "^(?!_)(?!.*?_$)[a-zA-Z0-9_\u4e00-\u9fa5]{2,8}$"
)

func init() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/userinfo", userinfo)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/register", register)
	// http.HandleFunc("/testjson", testjson)
	http.HandleFunc("/serverlist", serverlist)
	http.HandleFunc("/addserver", addserver)
	http.HandleFunc("/updateserver", updateserver)
	http.HandleFunc("/addplayer", addplayer)
	http.HandleFunc("/checkplayer", checkplayer)
	// http.ListenAndServeTLS(":8080", "853lab.com.cer", "853lab.com.key", nil)
}

func main() {
	log.Println("Hello, World!")
	log.Println(strconv.FormatInt(time.Now().Unix(), 10))
	log.Fatal(http.ListenAndServe(":8080", nil))
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// cfg := &tls.Config{
	// 	CipherSuites: []uint16{
	// 		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	// 		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	// 	},
	// 	Certificates:             []tls.Certificate{cert},
	// 	PreferServerCipherSuites: true,
	// 	InsecureSkipVerify:       true,
	// 	MinVersion:               tls.VersionTLS11,
	// 	MaxVersion:               tls.VersionTLS11,
	// }
	// log.Fatal(http.ListenAndServeTLS(":8080", "853lab.com.cer", "853lab.com.key", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//del
		user, _ := session.GetSession(w, r).GetAttr("user")

		t, err := template.ParseFiles("html/index.html")
		checkError(err)

		err = t.Execute(w, user)
		checkError(err)
		return
	}
	id := r.FormValue("id")
	key := r.FormValue("loginkey")
	if isEmpty(id, key) {
		message(w, r, 5001)
		return
	}
	isLogin, u := checklogin(id, key)
	if !isLogin || u == nil {
		message(w, r, 1006)
		return
	}
	w.Write([]byte(key))
}

// func testjson(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	con, _ := ioutil.ReadAll(r.Body)
// 	log.Println(string(con))
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(con)
// }

func serverlist(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		message(w, r, 5001)
		return
	}
	id := r.FormValue("id")
	key := r.FormValue("loginkey")
	if isEmpty(id, key) {
		message(w, r, 5001)
		return
	}
	isLogin, u := checklogin(id, key)
	if !isLogin || u == nil {
		message(w, r, 1006)
		return
	}
	lists := dao.ListServer()
	ret, _ := json.Marshal(lists)
	log.Println(string(ret))
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

func addserver(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		message(w, r, 5001)
		return
	}
	id := r.FormValue("id")
	key := r.FormValue("loginkey")
	ip := r.FormValue("ip")
	if isEmpty(id, key, ip) {
		message(w, r, 5001)
		return
	}
	if !validate(ip, "IP") {
		message(w, r, 5001)
		return
	}
	isLogin, u := checklogin(id, key)
	if !isLogin || u == nil {
		message(w, r, 1006)
		return
	}
	idi, err := strconv.Atoi(id)
	checkError(err)
	var name string
	if u.Name == "NULL" {
		name = u.Username
	} else {
		name = u.Name
	}
	timeinow := strconv.FormatInt(time.Now().Unix(), 10)
	checkError(err)
	server := dao.FindServerByID(id)
	if server == nil {
		server = &game.Server{
			ID:   idi,
			Name: name,
			IP:   ip,
			Time: timeinow,
		}
		server.SetServerKey(u)
		server.KickPlayers()
		dao.AddServer(server)
	} else {
		server.Name = name
		server.IP = ip
		server.Time = timeinow
		server.SetServerKey(u)
		server.KickPlayers()
		dao.UpdateServer(server)
	}
	log.Println(server.Key)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(server.Key))
}

func updateserver(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		message(w, r, 5001)
		return
	}
	id := r.FormValue("id")
	key := r.FormValue("loginkey")
	skey := r.FormValue("serverkey")
	if isEmpty(id, key, skey) {
		message(w, r, 5001)
		return
	}
	if !validate(skey, "Mkey") {
		message(w, r, 5001)
		return
	}
	isLogin, u := checklogin(id, key)
	if !isLogin || u == nil {
		message(w, r, 1006)
		return
	}
	server := dao.FindServerByID(id)
	if server == nil {
		message(w, r, 2005)
		return
	}
	if server.Key != skey {
		message(w, r, 2006)
		return
	}
	server.Time = strconv.FormatInt(time.Now().Unix(), 10)
	dao.UpdateServer(server)
	w.Write([]byte("true"))
}

func addplayer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		message(w, r, 5001)
		return
	}
	id := r.FormValue("id")
	key := r.FormValue("loginkey")
	sid := r.FormValue("serverid")
	if isEmpty(id, key, sid) {
		message(w, r, 5001)
		return
	}
	if !validate(sid, "ID") {
		message(w, r, 5001)
		return
	}
	isLogin, u := checklogin(id, key)
	if !isLogin || u == nil {
		message(w, r, 1006)
		return
	}
	server := dao.FindServerByID(sid)
	if server == nil {
		message(w, r, 2001)
		return
	}
	h := md5.New()
	h.Write([]byte(id + server.Key))
	muk := hex.EncodeToString(h.Sum(nil))[8:24]
	userjoin := dao.FindUserJoinByID(id)
	if userjoin == nil {
		idi, err := strconv.Atoi(id)
		checkError(err)
		userjoin = &game.UserJoin{
			ID:     idi,
			Server: 0,
			Key:    "NULL",
		}
		dao.AddUserJoin(userjoin)
	}
	if server.P1 == muk || server.P2 == muk || server.P3 == muk || server.P4 == muk {
		userjoin.Server = server.ID
		userjoin.Key = muk
	} else {
		if server.P1 == "NULL" {
			server.P1 = muk
		} else if server.P2 == "NULL" {
			server.P2 = muk
		} else if server.P3 == "NULL" {
			server.P3 = muk
		} else if server.P4 == "NULL" {
			server.P4 = muk
		} else {
			message(w, r, 2002)
			return
		}
		dao.UpdateServer(server)
		userjoin.Server = server.ID
		userjoin.Key = muk
	}
	dao.UpdateUserJoin(userjoin)
	joinserver := &game.JoinServer{
		IP:  server.IP,
		Key: userjoin.Key,
	}
	ret, _ := json.Marshal(joinserver)
	log.Println(string(ret))
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

func checkplayer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		message(w, r, 5001)
		return
	}
	id := r.FormValue("serverid")
	skey := r.FormValue("serverkey")
	pkey := r.FormValue("playerkey")
	if isEmpty(id, skey, pkey) {
		message(w, r, 5001)
		return
	}
	if !validate(skey, "Mkey") {
		message(w, r, 5001)
		return
	}
	if !validate(pkey, "PKey") {
		message(w, r, 5001)
		return
	}
	server := dao.FindServerByID(id)
	if server == nil {
		message(w, r, 2005)
		return
	}
	if server.Key != skey {
		message(w, r, 2006)
		return
	}
	pkeyList := strings.Split(pkey, ",")
	server.Time = strconv.FormatInt(time.Now().Unix(), 10)
	cp := ""
	for pki := 0; pki < len(pkeyList); pki++ {
		if server.P1 != pkeyList[pki] && server.P2 != pkeyList[pki] && server.P3 != pkeyList[pki] && server.P4 != pkeyList[pki] {
			log.Println("FAKE Key:" + pkeyList[pki])
		} else {
			if cp != "" {
				cp = cp + "," + pkeyList[pki]
			} else {
				cp = pkeyList[pki]
			}
		}
	}
	pkeyList = strings.Split(cp, ",")
	for pki := 0; pki < len(pkeyList); pki++ {
		if server.P1 == pkeyList[pki] {
			break
		} else if pki == 3 {
			server.P1 = "NULL"
		}
	}
	for pki := 0; pki < len(pkeyList); pki++ {
		if server.P2 == pkeyList[pki] {
			break
		} else if pki == 3 {
			server.P2 = "NULL"
		}
	}
	for pki := 0; pki < len(pkeyList); pki++ {
		if server.P3 == pkeyList[pki] {
			break
		} else if pki == 3 {
			server.P3 = "NULL"
		}
	}
	for pki := 0; pki < len(pkeyList); pki++ {
		if server.P4 == pkeyList[pki] {
			break
		} else if pki == 3 {
			server.P4 = "NULL"
		}
	}
	dao.UpdateServer(server)
	w.Write([]byte(cp))
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		loginHTML, err := ioutil.ReadFile("html/login.html")
		checkError(err)
		w.Write(loginHTML)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	encrypt := r.FormValue("encrypt")
	log.Println("login", username, password, encrypt)
	if isEmpty(username, password) {
		message(w, r, 1108)
		return
	}
	if encrypt != "yes" {
		if !validate(password, "Password") {
			message(w, r, 1001)
			return
		}
	} else {
		if !validate(password, "PasswordEnc") {
			message(w, r, 1001)
			return
		}
	}
	// log.Println("Passed validate")
	user := dao.FindUserByUsername(username)
	// if validate(username, "Username") {
	// 	user = dao.FindUserByUsername(username)
	// } else if validate(username, "Phone") {
	// 	user = dao.FindUserByPhone(username)
	// }
	if user == nil {
		message(w, r, 1001)
		return
	}
	if encrypt != "yes" {
		password = user.GetSaltedPassword(password, user.Salt)
	}
	if password != user.Password {
		message(w, r, 1001)
		return
	}
	// 登陆成功
	user.LoginKey = user.SetLoginKey(user.Password)
	dao.UpdateLoginKey(user)

	userinfo := &game.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Password: user.Password,
		LoginKey: user.LoginKey,
	}
	userinfo.SetPhone(user.Phone)
	sess := session.GetSession(w, r)
	sess.SetAttr("user", userinfo)
	ret, _ := json.Marshal(userinfo)
	log.Println(string(ret))
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

func userinfo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method == "GET" {
		sess := session.GetSession(w, r)
		user, exist := sess.GetAttr("user")
		if !exist {
			http.Redirect(w, r, "/", 302)
			return
		}
		ret, _ := json.Marshal(user)
		log.Println(string(ret))
		w.Header().Set("Content-Type", "application/json")
		w.Write(ret)
		return
	}

	con, _ := ioutil.ReadAll(r.Body)
	log.Println(string(con))

	cu := &game.UserInfo{}

	err := json.Unmarshal(con, &cu)
	if err != nil {
		message(w, r, 5001)
		return
	}

	// dec := json.NewDecoder(r.Body)
	// dec.DisallowUnknownFields()

	// err := dec.Decode(&cu)
	// if err != nil {
	// 	message(w, r, 5001)
	// 	return
	// }

	if !validate(strconv.Itoa(cu.ID), "ID") {
		// log.Println("ID")
		message(w, r, 5001)
		return
	}
	if !validate(cu.Username, "Username") {
		// log.Println("Username")
		message(w, r, 5001)
		return
	}
	if !validate(cu.Name, "Name") {
		// log.Println("Name")
		message(w, r, 5001)
		return
	}
	// if !validate(cu.Phone, "Phone") {
	// 	log.Println("Phone")
	// 	message(w, r, 5001)
	// 	return
	// }
	encrypt := true
	if !validate(cu.Password, "PasswordEnc") {
		if !validate(cu.Password, "Password") {
			// log.Println("Password")
			message(w, r, 5001)
			return
		}
		encrypt = false
	}
	if !validate(cu.LoginKey, "PasswordEnc") {
		// log.Println("PasswordEnc")
		message(w, r, 5001)
		return
	}

	userCheck := dao.FindUserByID(strconv.Itoa(cu.ID))
	if userCheck != nil {
		if cu.LoginKey != userCheck.LoginKey {
			message(w, r, 1006)
			return
		}
	} else {
		message(w, r, 1006)
		return
	}

	user1 := &game.User{
		ID:       userCheck.ID,
		Username: userCheck.Username,
		Name:     cu.Name,
		Phone:    userCheck.Phone,
		Password: userCheck.Password,
		LoginKey: userCheck.LoginKey,
		Salt:     userCheck.Salt,
	}

	if !encrypt {
		user1.SetPassword(cu.Password)
	}

	dao.UpdateUser(user1)

	userinfo := &game.UserInfo{
		ID:       user1.ID,
		Username: user1.Username,
		Name:     user1.Name,
		Password: user1.Password,
		LoginKey: user1.LoginKey,
	}
	userinfo.SetPhone(user1.Phone)

	ret, _ := json.Marshal(userinfo)
	log.Println(string(ret))
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
	return

	// POST 更新用户信息
	// username := r.FormValue("username")
	// id := r.FormValue("id")
	// loginkey := r.FormValue("loginkey")
	// name := r.FormValue("name")
	// password := r.FormValue("password")
	// phone := r.FormValue("phone")

	// if isEmpty(username, id, loginkey) {
	// 	message(w, r, 1108)
	// 	return
	// }

	// if !validate(id, "ID") {
	// 	message(w, r, 1102)
	// 	return
	// }

	// if !validate(loginkey, "PasswordEnc") {
	// 	message(w, r, 1102)
	// 	return
	// }

	// if !validate(username, "UserName") {
	// 	message(w, r, 1102)
	// 	return
	// }

	// if name != "" {
	// 	if !validate(name, "Name") {
	// 		message(w, r, 1109)
	// 		return
	// 	}
	// }

	// if password != "" {
	// 	if !validate(password, "Password") {
	// 		message(w, r, 1103)
	// 		return
	// 	}
	// }

	// if phone != "" {
	// 	if !validate(phone, "Phone") {
	// 		message(w, r, 1104)
	// 		return
	// 	}
	// }

	// switch user := user.(type) {
	// case *game.User:
	// 	if name != "" {
	// 		user.Name = name
	// 	}
	// 	if phone != "" {
	// 		user.Phone = phone
	// 	}
	// 	if password != "" {
	// 		user.Password = user.SetPassword(password)
	// 		user.LoginKey = user.SetLoginKey(user.Password)
	// 		sess.DelAttr("user")
	// 	}
	// 	dao.UpdateUser(user)
	// default:
	// 	log.Println(":userinfo:user.(type)", user)
	// }
	// http.Redirect(w, r, "/userinfo", 302)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		registerHTML, err := ioutil.ReadFile("html/register.html")
		checkError(err)
		_, err = w.Write(registerHTML)
		checkError(err)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	password2 := r.FormValue("password2")
	phone := r.FormValue("phone")

	if isEmpty(username, password, password2, phone) {
		message(w, r, 1108)
		return
	}

	if !validate(username, "Username") {
		message(w, r, 1102)
		return
	}

	if !validate(phone, "Phone") {
		message(w, r, 1104)
		return
	}

	userCheck := dao.FindUserByUsername(username)
	if userCheck != nil {
		message(w, r, 1105)
		return
	}
	userCheck = dao.FindUserByPhone(phone)
	if userCheck != nil {
		message(w, r, 1106)
		return
	}

	if !validate(password, "Password") {
		message(w, r, 1103)
		return
	}

	if !validate(password2, "Password") {
		message(w, r, 1103)
		return
	}

	if password != password2 {
		message(w, r, 1107)
		return
	}
	user := &game.User{
		Username: username,
		Name:     "NULL",
		Phone:    phone,
	}
	user.SetPassword(password)
	dao.AddUser(user)
	// message(w, r, 1100)
	errcode := &game.ErrCode{}
	errcode.SetCode(1100)
	ret, _ := json.Marshal(errcode)
	log.Println(string(ret))
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

func logout(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(w, r)
	sess.DelAttr("user")
	http.Redirect(w, r, "/", 302)
}

func message(w http.ResponseWriter, r *http.Request, code int) {
	errcode := &game.ErrCode{}
	errcode.SetCode(code)
	ret, _ := json.Marshal(errcode)
	log.Println(string(ret))
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func checklogin(id string, key string) (isLogin bool, u *game.User) {
	// if isEmpty(id, key) {
	// 	isLogin = false
	// 	return
	// }
	if !validate(id, "ID") {
		isLogin = false
		return
	}
	if !validate(key, "Mkey") {
		isLogin = false
		return
	}
	u = dao.FindUserByID(id)
	var muk string
	if u != nil {
		h := md5.New()
		h.Write([]byte(u.LoginKey))
		muk = hex.EncodeToString(h.Sum(nil))[8:24]
		if key != muk {
			isLogin = false
			return
		}
	} else {
		isLogin = false
		return
	}
	isLogin = true
	return
}

func isEmpty(strs ...string) (isEmpty bool) {
	for _, str := range strs {
		str = strings.TrimSpace(str)
		if str == "" || len(str) == 0 {
			isEmpty = true
			return
		}
	}
	isEmpty = false
	return
}

func validate(userinput string, what string) bool {
	// reg := regexp.MustCompile("")
	switch what {
	case "Username":
		reg := regexp.MustCompile(regularUsername)
		return reg.MatchString(userinput)
	case "Password":
		reg := regexp.MustCompile(regularPassword)
		return reg.MatchString(userinput)
	case "PasswordEnc":
		reg := regexp.MustCompile(regularPasswordEnc)
		return reg.MatchString(userinput)
	case "Phone":
		reg := regexp.MustCompile(regularPhone)
		return reg.MatchString(userinput)
	case "Name":
		reg := regexp.MustCompile(regularName)
		return reg.MatchString(userinput)
	case "ID":
		reg := regexp.MustCompile(regularID)
		return reg.MatchString(userinput)
	case "Mkey":
		reg := regexp.MustCompile(regularMkey)
		return reg.MatchString(userinput)
	case "IP":
		reg := regexp.MustCompile(regularIP)
		return reg.MatchString(userinput)
	case "PKey":
		reg := regexp.MustCompile(regularPkey)
		return reg.MatchString(userinput)
	}
	return false
}
