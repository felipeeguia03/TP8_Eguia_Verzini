package clients

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"backend/dao"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBClient *gorm.DB
)

// Se mantiene init() para NO romper tu ejecución local.
// Lee env si existen; si no, usa tus defaults locales.
// Reintenta algunas veces para Render (DB tarda).
func init() {
	// Debug: mostrar qué variables están disponibles
	log.Info("🔍 Checking database configuration...")
	databaseURL := getEnv("DATABASE_URL", "")
	if databaseURL != "" {
		log.Info("✅ DATABASE_URL found (length: " + fmt.Sprintf("%d", len(databaseURL)) + ")")
	} else {
		log.Warn("⚠️ DATABASE_URL not found")
		log.Info("Checking individual DB variables...")
		log.WithFields(log.Fields{
			"DB_HOST": getEnv("DB_HOST", "NOT SET"),
			"DB_USER": getEnv("DB_USER", "NOT SET"),
			"DB_NAME": getEnv("DB_NAME", "NOT SET"),
			"DB_PORT": getEnv("DB_PORT", "NOT SET"),
			"has_password": getEnv("DB_PASSWORD", "") != "",
		}).Info("Individual variables status")
	}
	
	var dsn string
	
	// Si existe DATABASE_URL (Railway, Render, etc.), usarlo directamente
	if databaseURL != "" {
		dsn = databaseURL
		log.Info("✅ Using DATABASE_URL for connection")
	} else {
		log.Warn("⚠️ DATABASE_URL not found, using individual variables")
		// Fallback a variables individuales (compatibilidad con Render, local, etc.)
		dbName := getEnv("DB_NAME", "final_clj4")
		dbUser := getEnv("DB_USER", "admin")
		dbPassword := getEnv("DB_PASSWORD", "")
		dbHost := getEnv("DB_HOST", "127.0.0.1")
		portStr := getEnv("DB_PORT", "5432")
		dbPort, err := strconv.Atoi(portStr)
		if err != nil {
			dbPort = 5432
		}

		// SSL mode para PostgreSQL
		sslMode := getEnv("DB_SSLMODE", "require")
		if sslMode == "" {
			sslMode = "require"
		}

		// Log de configuración (sin password) - mostrar qué variables se están usando
		log.WithFields(log.Fields{
			"host": dbHost,
			"port": dbPort,
			"user": dbUser,
			"db":   dbName,
			"ssl":  sslMode,
			"has_password": dbPassword != "",
		}).Info("Using individual DB variables for connection")
		
		// Validar que no estemos usando defaults en producción (Render, Railway, etc.)
		// Si PORT está configurado (típico de Render), no deberíamos usar localhost
		port := getEnv("PORT", "")
		if port != "" && (dbHost == "127.0.0.1" || dbHost == "localhost") {
			log.Error("❌ ERROR: Using localhost in production environment!")
			log.Error("Current configuration:")
			log.Errorf("  DB_HOST=%s (should be Railway host)", dbHost)
			log.Errorf("  DB_USER=%s (should be Railway user)", dbUser)
			log.Errorf("  DB_NAME=%s", dbName)
			log.Errorf("  DB_PORT=%d", dbPort)
			log.Error("Please configure DATABASE_URL or all DB_* variables in Render Environment Variables")
			panic("Database configuration error: localhost detected in production. Set DATABASE_URL or individual DB_* variables.")
		}

		// Validar que tengamos password si no es localhost
		if dbHost != "127.0.0.1" && dbHost != "localhost" && dbPassword == "" {
			log.Error("❌ ERROR: DB_PASSWORD is required for remote database!")
			panic("Database configuration error: DB_PASSWORD is missing for remote host.")
		}

		// Formato DSN para PostgreSQL: host=host user=user password=password dbname=dbname port=port sslmode=mode
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			dbHost, dbUser, dbPassword, dbName, dbPort, sslMode)
	}

	const maxRetries = 10
	var lastErr error

	for i := 1; i <= maxRetries; i++ {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			DBClient = db
			// Extraer info de conexión para el log
			dbHost := getEnv("DB_HOST", getEnv("DATABASE_URL", "unknown"))
			dbName := getEnv("DB_NAME", getEnv("PGDATABASE", "unknown"))
			log.WithFields(log.Fields{"host": dbHost, "db": dbName, "attempt": i}).Info("DB connected")
			return
		}
		lastErr = err
		log.WithFields(log.Fields{"attempt": i, "err": err}).Warn("DB connection failed, retrying...")
		time.Sleep(3 * time.Second)
	}

	panic(fmt.Errorf("error connecting to DB after retries: %v", lastErr))
}

func StartDB() {
	var (
		user         dao.User
		course       dao.Course
		subscription dao.Subscription
		comment      dao.Comment
		file         dao.File
	)

	if err := DBClient.AutoMigrate(&user, &course, &subscription, &comment, &file); err != nil {
		panic(fmt.Errorf("error creating entities: %v", err))
	}

	log.Info("Database initialized")
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func CreateUser(NewUser dao.User) error {
	var user dao.User

	result := DBClient.Where("Email = ?", NewUser.Email).First(&user)

	if result.Error == nil {
		return fmt.Errorf("user with email %s already exists", NewUser.Email)
	}

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return fmt.Errorf("error checking for existing user for email %s", NewUser.Email)
	}

	result = DBClient.Where("Nickname = ?", NewUser.Nickname).First(&user)

	if result.Error == nil {
		return fmt.Errorf("user with nickname %s already exists", NewUser.Nickname)
	}

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return fmt.Errorf("error checking for existing user for nickname %s", NewUser.Nickname)
	}

	result = DBClient.Create(&NewUser)
	if result.Error != nil {
		return fmt.Errorf("error creating user for nickname %s and email %s", NewUser.Nickname, NewUser.Email)
	}

	log.Debug("User created: ", NewUser)
	return nil

}

func GetUserById(ID int64) (dao.User, error) {
	var user dao.User

	result := DBClient.Where("id = ?", ID).First(&user)

	if result.Error != nil {
		return user, fmt.Errorf("not found user with ID: %d", ID)
	}

	log.Debug("User: ", user)
	return user, nil
}

func GetCourses() ([]dao.Course, error) {
	var courses []dao.Course

	result := DBClient.Find(&courses)

	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving courses: %s", result.Error)
	}

	return courses, nil
}

func GetUserByEmail(email string) (dao.User, error) {
	var user dao.User

	result := DBClient.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return user, fmt.Errorf("not found user with email: %s", email)
	}

	log.Debug("User: ", user)
	return user, nil
}

func GetCourseById(ID int64) (dao.Course, error) {
	var course dao.Course

	result := DBClient.Where("id = ?", ID).First(&course)

	if result.Error != nil {
		return course, fmt.Errorf("not found course with ID: %d", ID)
	}

	log.Debug("Course: ", course)
	return course, nil
}

func GetCoursewithQuery(query string) ([]dao.Course, error) {
	var courses []dao.Course

	result := DBClient.Where("title LIKE ? OR description LIKE ? OR category LIKE ? OR requirement LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%").Find(&courses)

	if result.Error != nil {
		return nil, fmt.Errorf("not found course with filter: %s", query)
	}

	return courses, nil
}

func GetCourseIdsByUserId(userID int64) ([]int64, error) {
	var subscriptions []dao.Subscription
	var courseIDs []int64

	result := DBClient.Where("user_id = ?", userID).Find(&subscriptions)
	if result.Error != nil {
		return nil, fmt.Errorf("error finding subscriptions for user ID: %d, %v", userID, result.Error)
	}

	for _, subscription := range subscriptions {
		courseIDs = append(courseIDs, subscription.Course_Id)
	}

	return courseIDs, nil
}

func InsertSubscription(userID int64, courseID int64) error {
	var subscription dao.Subscription

	result := DBClient.Where("user_id = ? AND course_id = ?", userID, courseID).First(&subscription)

	if result.Error == nil {
		return fmt.Errorf("subscription already exists for user %d and course %d", userID, courseID)
	}

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return fmt.Errorf("error checking for existing subscription for user %d and course %d", userID, courseID)
	}

	NewSubscription := dao.Subscription{
		User_Id:   userID,
		Course_Id: courseID,
	}

	result = DBClient.Create(&NewSubscription)
	if result.Error != nil {
		return fmt.Errorf("error inserting subscription for user %d and course %d", userID, courseID)
	}

	log.Debug("Subscription created: ", NewSubscription)
	return nil

}

func DeleteCourseById(courseID int64) error {
	var course dao.Course

	result := DBClient.Where("id = ?", courseID).First(&course)

	if result.Error != nil {
		return fmt.Errorf("not found course with ID: %d", courseID)
	}

	result = DBClient.Delete(&course)
	if result.Error != nil {
		return fmt.Errorf("error deleting course for ID: %d", courseID)
	}

	log.Debug("Course Deleted: ", course)
	return nil
}

func DeleteSubscriptionById(courseID int64) error {
	var subscription dao.Subscription

	result := DBClient.Where("course_id = ?", courseID).Delete(&subscription)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return fmt.Errorf("error deleting course with id %d in subscription table", courseID)
	}

	return nil
}

func CreateCourse(NewCourse dao.Course) error {
	var course dao.Course

	result := DBClient.Where("Title = ?", NewCourse.Title).First(&course)

	if result.Error == nil {
		return fmt.Errorf("course with title %s already exists", NewCourse.Title)
	}

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return fmt.Errorf("error checking for existing course for title %s", NewCourse.Title)
	}

	result = DBClient.Create(&NewCourse)
	if result.Error != nil {
		return fmt.Errorf("error creating course for title: %s", NewCourse.Title)
	}

	log.Debug("Course created: ", NewCourse)
	return nil
}

func UpdateCourse(courseID int64, courseUpdate dao.Course) error {
	var course dao.Course

	result := DBClient.Where("id = ?", courseID).First(&course)
	if result.Error != nil {
		return fmt.Errorf("not found course with ID: %d", courseID)
	}

	course.Title = courseUpdate.Title
	course.Description = courseUpdate.Description
	course.Category = courseUpdate.Category
	course.Duration = courseUpdate.Duration
	course.Instructor = courseUpdate.Instructor
	course.Requirement = courseUpdate.Requirement

	result = DBClient.Save(&course)
	if result.Error != nil {
		return fmt.Errorf("error updating course with ID:  %d", course.Id)
	}

	log.Debug("Course updated: ", course)
	return nil
}

func InsertComment(userID int64, courseID int64, text string) error {
	var comment dao.Comment
	result := DBClient.Where("user_id = ? AND course_id = ? AND text = ?", userID, courseID, text).First(&comment)

	if result.Error == nil {
		return fmt.Errorf("comment already exists for user %d and course %d", userID, courseID)
	}

	NewComment := dao.Comment{
		User_Id:   userID,
		Course_Id: courseID,
		Text:      text,
	}

	result = DBClient.Create(&NewComment)
	if result.Error != nil {
		return fmt.Errorf("error inserting comment for user %d and course %d", userID, courseID)
	}

	log.Debug("Comment created: ", NewComment)
	return nil
}

func GetCommentById(commentID int64) (dao.Comment, error) {
	var comment dao.Comment

	result := DBClient.Where("id = ?", commentID).First(&comment)

	if result.Error != nil {
		return comment, fmt.Errorf("not found comment with ID: %d", commentID)
	}

	log.Debug("Comment: ", comment)
	return comment, nil
}

func GetCommentsByCourseId(courseID int64) ([]int64, error) {
	var comments []dao.Comment
	var commentIDs []int64

	result := DBClient.Where("course_Id = ?", courseID).Find(&comments)
	if result.Error != nil {
		return nil, fmt.Errorf("error finding comments for course ID: %d, %v", courseID, result.Error)
	}

	for _, comment := range comments {
		commentIDs = append(commentIDs, comment.Id)
	}

	return commentIDs, nil
}

func SaveFile(NewFile dao.File) error {
	var course dao.Course
	var user dao.User

	result := DBClient.Where("id = ?", NewFile.Course_Id).First(&course)
	if result.Error != nil {
		return fmt.Errorf("not found course with ID: %d", NewFile.Course_Id)
	}

	result = DBClient.Where("id = ?", NewFile.User_Id).First(&user)
	if result.Error != nil {
		return fmt.Errorf("not found user with ID: %d", NewFile.User_Id)
	}

	result = DBClient.Create(&NewFile)
	if result.Error != nil {
		return fmt.Errorf("error creating file record: %v", result.Error)
	}
	log.Debug("File record created: ", NewFile)
	return nil
}
