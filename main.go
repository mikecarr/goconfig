package main

import (
	"encoding/json"
	"image/color"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/ssh"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Settings structure to hold the configuration
type Settings struct {
	IP       string `json:"ip"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoadSettings reads settings from a JSON file
func LoadSettings(filename string) (Settings, error) {
	var settings Settings

	// Read the file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return settings, err
	}

	// Unmarshal JSON into settings struct
	err = json.Unmarshal(data, &settings)
	return settings, err
}

// ConnectToSSH establishes an SSH connection using the provided credentials
func ConnectToSSH(ip, username, password string) error {
	// Configure the SSH client
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Not secure for production
	}

	// Establish the SSH connection
	conn, err := ssh.Dial("tcp", ip+":22", config) // Default SSH port is 22
	if err != nil {
		return err
	}
	defer conn.Close() // Close connection when done

	return nil
}

// UpdateLog appends a log message to the log entry
func UpdateLog(logEntry *widget.Entry, message string) {
	logEntry.SetText(logEntry.Text + "\n" + message)
}

// SaveLog saves the log entry contents to a specified file
func SaveLog(logEntry *widget.Entry, window fyne.Window) {
	// Open a file dialog to choose where to save the log
	dialog.ShowFileSave(func(file fyne.URIWriteCloser, err error) {
		if err != nil || file == nil {
			UpdateLog(logEntry, "File save canceled or error occurred.")
			return
		}
		// Write log entry contents to the file
		if _, err := file.Write([]byte(logEntry.Text)); err != nil {
			UpdateLog(logEntry, "Error saving log: "+err.Error())
		} else {
			UpdateLog(logEntry, "Log saved successfully.")
		}
	}, window)
}

func main() {
	// Create a new application with a unique ID
	myApp := app.NewWithID("com.example.openipc") // Use a unique ID for your application
	myWindow := myApp.NewWindow("Tab Control Example")

	// Set the initial size of the window
	myWindow.Resize(fyne.NewSize(800, 600)) // Width: 800, Height: 600

	// Load settings from file
	settings, err := LoadSettings("settings.json")
	if err != nil {
		log.Fatalf("Failed to load settings: %v", err)
	}

	// Create labels and entry boxes for IP, Username, and Password
	ipLabel := widget.NewLabel("IP:")
	ipEntry := widget.NewEntry()
	ipEntry.SetText(settings.IP) // Set default value for IP

	usernameLabel := widget.NewLabel("Username:")
	usernameEntry := widget.NewEntry()
	usernameEntry.SetText(settings.Username) // Set default value for Username

	passwordLabel := widget.NewLabel("Password:")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetText(settings.Password) // Set default value for Password

	// Create a label and entry for 5.8 Frequency
	frequencyLabel := widget.NewLabel("5.8 Frequency:")
	frequencyEntry := widget.NewEntry()

	// Create a log entry field
	logEntry := widget.NewMultiLineEntry()
	logEntry.SetPlaceHolder("Logs will appear here...") // Placeholder for the log entry
	logEntry.SetText("Log output:\n")                    // Initial log output

	// Create a "Connect" button and set its size
	connectButton := widget.NewButton("Connect", func() {
		// Handle connect action here
		ip := ipEntry.Text
		username := usernameEntry.Text
		password := passwordEntry.Text

		// Attempt to connect via SSH
		err := ConnectToSSH(ip, username, password)
		if err != nil {
			log.Println("Connection failed:", err)
			UpdateLog(logEntry, "Connection failed: "+err.Error())
		} else {
			log.Println("Connected successfully to:", ip, "with username:", username)
			UpdateLog(logEntry, "Connected successfully to: "+ip+" with username: "+username)
		}
	})

	// Create a "Save Log" button
	saveLogButton := widget.NewButton("Save Log", func() {
		SaveLog(logEntry, myWindow) // Pass the window to the SaveLog function
	})

	// Create content for Tab 1 with the frequency entry aligned horizontally
	tab1Content := container.NewVBox(
		widget.NewLabel("Content of Tab 1"),
		container.NewHBox(frequencyLabel, frequencyEntry), // Align label and entry horizontally
	)

	// Create a TabContainer
	tabs := container.NewAppTabs(
		container.NewTabItem("Tab 1", tab1Content),
		container.NewTabItem("Tab 2", container.NewVBox(widget.NewLabel("Content of Tab 2"))),
		container.NewTabItem("Tab 3", container.NewVBox(widget.NewLabel("Content of Tab 3"))),
	)

	// Create a border rectangle with a color
	border := canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 0, A: 255}) // Black color for the border
	border.SetMinSize(fyne.NewSize(800, 300))                        // Set minimum size to fit the tabs

	// Create a container for the border and the tabs
	tabContainer := container.NewMax(border, tabs)

	// Vertical layout for labels and entries
	formContainer := container.NewVBox(
		container.NewHBox(ipLabel, ipEntry),
		container.NewHBox(usernameLabel, usernameEntry),
		container.NewHBox(passwordLabel, passwordEntry),
		connectButton,
		saveLogButton,
	)

	// Wrap the log entry in a scroll container
	scrollLogEntry := container.NewScroll(logEntry)

	// Create a vertical container to hold the bordered tabs, form, and log entry
	content := container.NewVBox(tabContainer, formContainer, scrollLogEntry)

	// Set the content of the window
	myWindow.SetContent(content)

	// Show and run the application
	myWindow.ShowAndRun()
}
