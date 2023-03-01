package controllers

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/divrhino/fruitful-pdf/models"
	"github.com/gofiber/fiber/v2"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func CreateCertificate(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("missing file parameter")
	}
	result := models.DetailCertificate{
		Name:     c.FormValue("name"),
		LastName: c.FormValue("last_name"),
		Course:   c.FormValue("course"),
		Date:     c.FormValue("date"),
		Template: c.FormValue("template"),
	}

	m := pdf.NewMaroto(consts.Landscape, consts.A4)

	// err = m.FileImage("images/10.png", props.Rect{
	// 	Left: 30,
	// })
	// if err != nil {
	// 	fmt.Println("Image file was not loaded 😱 - ", err)
	// }

	//select tamplate
	if result.Template == "1" {
		TemplateFirst(m, result, file)
	} else if result.Template == "2" {
		TemplateSecond(m, result, file)
	}

	//Gen file name
	t := time.Now()
	ynmStr := t.Format("200601")

	pdfBytes, _ := m.Output()
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename="+ynmStr+fmt.Sprint(t.Day())+fmt.Sprint(t.Hour())+fmt.Sprint(t.Minute())+fmt.Sprint(t.Second())+".pdf")
	fmt.Println("PDF saved successfully")

	return c.Send(pdfBytes.Bytes())
}

func CreateCertificateCSV(c *fiber.Ctx) error {
	var fileType, fileName string
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("error file is null")
		return c.Status(401).JSON(err)
	} else {
		src, err := file.Open()
		if err != nil {
			fmt.Println("error can't open file")
			return c.Status(500).JSON(err)
		}
		fileByte, _ := ioutil.ReadAll(src)
		fileType = http.DetectContentType(fileByte)
		if fileType == "application/csv" {
			fileName = "uploads/" + strconv.FormatInt(time.Now().Unix(), 10) + ".csv"
		} else {
			fileName = "uploads/" + strconv.FormatInt(time.Now().Unix(), 10) + ".csv"
		}

		err = ioutil.WriteFile(fileName, fileByte, 0777)
		if err != nil {
			return c.Status(500).JSON(err)
		}
		defer src.Close()
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(err)
	}
	courses := []models.DetailCertificate{}
	df := csv.NewReader(src)
	data, _ := df.ReadAll()
	for _, v := range data[1:] {
		courses = append(courses,
			models.DetailCertificate{
				Name:     v[0],
				LastName: v[1],
				Course:   v[2],
				Date:     v[3],
				Template: v[4]})
	}

	// t := time.Now()
	// ynmStr := t.Format("200601")
	for i := 0; i < len(courses); i++ {
		m := pdf.NewMaroto(consts.Landscape, consts.A4)
		TemplateTest(m, courses[i])

		//out put return เป็น base 64
		pdfBytes, _ := m.Output()
		//len bytes มี 5 ตัว
		pdfSlice := pdfBytes.Bytes()[:]
		//len ของ slice มี 5 ตัว

		// Create a new buffer to write the zip archive to
		var buf bytes.Buffer
		// Create a new zip archive
		zipWriter := zip.NewWriter(&buf)

		// Create a new file header for the PDF file
		fileHeader := &zip.FileHeader{
			Name:   "example.pdf",
			Method: zip.Deflate,
		}
		// Write the file header to the zip archive
		pdfFileWriter, err := zipWriter.CreateHeader(fileHeader)
		if err != nil {
			fmt.Println(err)

		}
		// Write the PDF file content to the zip archive
		_, err = pdfFileWriter.Write(pdfSlice)
		if err != nil {
			fmt.Println(err)

		}
		// Close the zip archive
		err = zipWriter.Close()
		if err != nil {
			fmt.Println(err)
		}
		// Write the zip archive to a file
		err = ioutil.WriteFile("example.zip", buf.Bytes(), 0644)
		if err != nil {
			fmt.Println(err)
		}
		//ออกเป็นคนสุดท้าย
	}

	return c.SendString("success")
}

func TemplateTest(m pdf.Maroto, data models.DetailCertificate) {
	// m.AddPage()
	m.RegisterHeader(func() {
		m.Row(40, func() {
			m.Col(6, func() {
				err := m.FileImage("images/header-logo.png", props.Rect{
					Left:    1,
					Percent: 100,
				})
				if err != nil {
					fmt.Println("Image file was not loaded 😱 - ", err)
				}
			})
		})
	})

	m.Row(20, func() {
		m.Col(12, func() {
			m.Text("CERTIFICATE", props.Text{
				Style:  consts.Bold,
				Align:  consts.Center,
				Family: consts.Arial,
				Size:   54,
				Color:  getDarkPurpleColor(),
			})
			m.Text("OF COMPLETION", props.Text{
				Top:   150,
				Style: consts.Bold,
				Align: consts.Center,
				Size:  32,
				Color: getDarkPurpleColor(),
			})
		})
	})
	m.Row(18, func() {
		m.Col(12, func() {
			m.Text("This certificate is proudly presented to", props.Text{
				Top:   120,
				Style: consts.Bold,
				Align: consts.Center,
				Size:  14,
				Color: getDarkPurpleColor(),
			})
		})
	})
	m.Row(20, func() {
		m.Col(12, func() {
			m.Text(data.Name+" "+data.LastName, props.Text{
				Top:   120,
				Style: consts.Bold,
				Align: consts.Center,
				Size:  26,
				Color: getDarkPurpleColor(),
			})
		})
	})

	m.Row(15, func() {
		m.Col(12, func() {
			m.Text("Successfully completed and received a passing grade in"+" "+data.Course, props.Text{
				Top:   120,
				Style: consts.Bold,
				Align: consts.Center,
				Size:  16,
				Color: getDarkPurpleColor(),
			})
		})
	})

	m.Row(20, func() {
		m.Col(6, func() {
			m.Text("SIGNATURE", props.Text{
				Left:  10,
				Top:   120,
				Style: consts.Normal,
				Align: consts.Left,
				Size:  18,
				Color: getDarkPurpleColor(),
			})

			err := m.FileImage("images/sign.png", props.Rect{
				Top:  30,
				Left: 10,
			})
			if err != nil {
				fmt.Println("Image file was not loaded 😱 - ", err)
			}
		})

		m.Col(6, func() {
			m.Text("DATE", props.Text{
				Top:   120,
				Style: consts.Normal,
				Align: consts.Center,
				Size:  18,
				Color: getDarkPurpleColor(),
			})
		})

	})

	m.Row(15, func() {
		m.Col(12, func() {
			m.Text(data.Date, props.Text{
				Top:   90,
				Left:  140,
				Style: consts.Normal,
				Align: consts.Center,
				Size:  18,
				Color: getDarkPurpleColor(),
			})
		})

	})
}

func TemplateFirst(m pdf.Maroto, data models.DetailCertificate, file *multipart.FileHeader) {
	m.RegisterHeader(func() {
		m.Row(40, func() {
			m.Col(6, func() {
				err := m.Base64Image(ImageUploadEndPoint(file), consts.Png, props.Rect{
					Left:    1,
					Percent: 100,
				})
				if err != nil {
					fmt.Println("Image file was not loaded 😱 - ", err)
				}
			})
		})
	})

	m.Row(20, func() {
		m.Col(12, func() {
			m.Text("CERTIFICATE", props.Text{
				Style:  consts.Bold,
				Align:  consts.Center,
				Family: consts.Arial,
				Size:   54,
				Color:  getDarkPurpleColor(),
			})
			m.Text("OF COMPLETION", props.Text{
				Top:   150,
				Style: consts.Bold,
				Align: consts.Center,
				Size:  32,
				Color: getDarkPurpleColor(),
			})
		})
	})
	m.Row(18, func() {
		m.Col(12, func() {
			m.Text("This certificate is proudly presented to", props.Text{
				Top:   120,
				Style: consts.Bold,
				Align: consts.Center,
				Size:  14,
				Color: getDarkPurpleColor(),
			})
		})
	})
	m.Row(20, func() {
		m.Col(12, func() {
			m.Text(data.Name+" "+data.LastName, props.Text{
				Top:   120,
				Style: consts.Bold,
				Align: consts.Center,
				Size:  26,
				Color: getDarkPurpleColor(),
			})
		})
	})

	m.Row(15, func() {
		m.Col(12, func() {
			m.Text("Successfully completed and received a passing grade in"+" "+data.Course, props.Text{
				Top:   120,
				Style: consts.Bold,
				Align: consts.Center,
				Size:  16,
				Color: getDarkPurpleColor(),
			})
		})
	})

	m.Row(20, func() {
		m.Col(6, func() {
			m.Text("SIGNATURE", props.Text{
				Left:  10,
				Top:   120,
				Style: consts.Normal,
				Align: consts.Left,
				Size:  18,
				Color: getDarkPurpleColor(),
			})

			err := m.FileImage("images/sign.png", props.Rect{
				Top:  30,
				Left: 10,
			})
			if err != nil {
				fmt.Println("Image file was not loaded 😱 - ", err)
			}
		})

		m.Col(6, func() {
			m.Text("DATE", props.Text{
				Top:   120,
				Style: consts.Normal,
				Align: consts.Center,
				Size:  18,
				Color: getDarkPurpleColor(),
			})
		})

	})

	m.Row(15, func() {
		m.Col(12, func() {
			m.Text(data.Date, props.Text{
				Top:   90,
				Left:  140,
				Style: consts.Normal,
				Align: consts.Center,
				Size:  18,
				Color: getDarkPurpleColor(),
			})
		})

	})
}

func TemplateSecond(m pdf.Maroto, data models.DetailCertificate, file *multipart.FileHeader) {
	m.RegisterHeader(func() {
		m.Row(40, func() {
			m.Col(6, func() {
				err := m.Base64Image(ImageUploadEndPoint(file), consts.Png, props.Rect{
					Left:    1,
					Percent: 100,
				})
				if err != nil {
					fmt.Println("Image file was not loaded 😱 - ", err)
				}
			})
		})
	})

	m.Row(20, func() {
		m.Col(12, func() {
			m.Text("CERTIFICATE", props.Text{
				Style:  consts.Bold,
				Align:  consts.Left,
				Family: consts.Arial,
				Size:   46,
				Color:  getDarkPurpleColor(),
			})
			m.Text("OF COMPLETION", props.Text{
				Top:   150,
				Style: consts.Bold,
				Align: consts.Left,
				Size:  32,
				Color: getDarkPurpleColor(),
			})
		})
	})
	m.Row(18, func() {
		m.Col(12, func() {
			m.Text("This certificate is proudly presented to", props.Text{
				Top:   120,
				Style: consts.Bold,
				Align: consts.Left,
				Size:  14,
				Color: getDarkPurpleColor(),
			})
		})
	})
	m.Row(20, func() {
		m.Col(12, func() {
			m.Text(data.Name+" "+data.LastName, props.Text{
				Top:   120,
				Style: consts.Bold,
				Align: consts.Left,
				Size:  26,
				Color: getDarkPurpleColor(),
			})
		})
	})

	m.Row(15, func() {
		m.Col(12, func() {
			m.Text("Successfully completed and received a passing grade in"+" "+data.Course, props.Text{
				Top:   120,
				Style: consts.Bold,
				Align: consts.Left,
				Size:  16,
				Color: getDarkPurpleColor(),
			})
		})
	})

	m.Row(20, func() {
		m.Col(6, func() {
			m.Text("SIGNATURE", props.Text{
				Left:  10,
				Top:   120,
				Style: consts.Normal,
				Align: consts.Left,
				Size:  18,
				Color: getDarkPurpleColor(),
			})

			err := m.FileImage("images/sign.png", props.Rect{
				Top:  30,
				Left: 10,
			})
			if err != nil {
				fmt.Println("Image file was not loaded 😱 - ", err)
			}
		})

		m.Col(6, func() {
			m.Text("DATE", props.Text{
				Top:   120,
				Style: consts.Normal,
				Align: consts.Left,
				Size:  18,
				Color: getDarkPurpleColor(),
			})
		})

	})

	m.Row(15, func() {
		m.Col(12, func() {
			m.Text(data.Date, props.Text{
				Top:   90,
				Left:  140,
				Style: consts.Normal,
				Align: consts.Left,
				Size:  18,
				Color: getDarkPurpleColor(),
			})
		})

	})
}

func getDarkPurpleColor() color.Color {
	return color.Color{
		Red:   88,
		Green: 80,
		Blue:  99,
	}
}

func ImageUploadEndPoint(file *multipart.FileHeader) string {
	fileData, err := file.Open()
	if err != nil {
		return ""
	}
	defer fileData.Close()

	imgData, err := ioutil.ReadAll(fileData)
	if err != nil {
		return ""
	}

	encodedImg := base64.StdEncoding.EncodeToString(imgData)

	return encodedImg
}

// upload file csv and read file
func UploadFileGo(c *fiber.Ctx) error {
	var personfile models.PersonUpload
	var fileType, fileName string
	// var fileSize int64
	isSuccess := true
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("error file is null")
		return c.Status(401).JSON(err)
	} else {
		src, err := file.Open()
		if err != nil {
			fmt.Println("error can't open file")
			return c.Status(500).JSON(err)
		}
		fileByte, _ := ioutil.ReadAll(src)
		fileType = http.DetectContentType(fileByte)
		if fileType == "application/csv" {
			fileName = "upload/" + strconv.FormatInt(time.Now().Unix(), 10) + ".csv"
		} else {
			fileName = "upload/" + strconv.FormatInt(time.Now().Unix(), 10) + ".csv"
		}

		err = ioutil.WriteFile(fileName, fileByte, 0777)
		if err != nil {
			return c.Status(500).JSON(err)
		}
		defer src.Close()
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(err)
	}
	courses := []models.DetailCertificate{}
	df := csv.NewReader(src)
	data, _ := df.ReadAll()
	for _, v := range data[2:] {
		courses = append(courses,
			models.DetailCertificate{
				Name:     v[0],
				LastName: v[1],
				Course:   v[2],
				Template: v[3],
				Date:     v[4]})
	}
	if isSuccess {
		personfile = models.PersonUpload{
			NameFile: "upload file success",
			File: struct {
				Data interface{}
			}{
				Data: courses,
			},
		}
	} else {
		personfile = models.PersonUpload{
			NameFile: "upload file success",
			File: struct {
				Data interface{}
			}{
				Data: courses,
			},
		}
	}
	return c.Status(200).JSON(personfile)
}
