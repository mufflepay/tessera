package controllers

import (
	"backend/models"
	"backend/repos"
	"backend/util"
	"backend/util/cloudStorage"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

// R2
const r2BucketName = "golangtest"
const accountId = "13c2e769e85f95a75dd19c3cf88db086"
const accessKeyId = "f7270102cf8c6abcd4cf65a811503d93"
const accessKeySecret = "9912d68a213f61dffc599adc4a7461aec61f2fa351c803f4169720503a4d2828"

// OCI
const tenancyOCID = "ocid1.tenancy.oc1..aaaaaaaapop6u3bytsheistqxkxyj7cmhetdqz7z3qkcwh7rei4x2rv5marq"

// const compartmentOCID = "ocid1.compartment.oc1..aaaaaaaack5um6vtucfivganmrrapf3csvtwkf53s3pjtqqxgnliag2buqba"
const userOCID = "ocid1.user.oc1..aaaaaaaaux5ouw6jmiabceva7fpyni4odvj4ccgnvdbfj3or4lkacnncs3va"
const region = "af-johannesburg-1"
const fingerprint = "88:b3:eb:96:9e:11:d5:fe:5d:5e:47:69:df:5f:12:56"
const namespace = "axrj9wzaeep0"
const ociBucketName = "bucket-20230812-0418"

const privateKey = "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDT6IER7gYQ+SJf\nDVsvQbsD8TUeu3Ae4wrHH/lUWNN61rLpuP4kprzE8sBJUE3r3uCSe9vNYSHR/d3p\nPmV750JTti2sMeLdLH/IQmee7nrdSzjLcwD7++1fezH2Cc4ts3ryNkcJCX9yI9Aa\ncTbtae0SCFfgig8V4ovo3bMfeznFXgsUCHrvXA1+U0hBMA8LQovbOoBgO+AsAJ9i\nIZwsCICM18tEEm9O9tu2K8oLG+hkYAbwioDu0Tc/axY5Dv2ba/+Z5ja0NUSfRJI1\n96Zsz+lYqRbkQDLFeggjQ0cc8xy62ULVGchjUYs06ovrfJLZFZ2EV51Bpe/H/UHW\nbY9gPWfNAgMBAAECggEAGkFcfJQ5PDEiCRTmj1xdh4eDRWOD6M/IrhNQGRtIWJfx\nYvQAyRm/mcdZ+7tvbCIZQQ0HltLKFfKWZCfO+yMUHLsdvZAQw8aXroBLxm1V28VO\nLCb7oDz+lC33arycRx1NxmTjOENs/v9g8WXFoHTXYryV5sQ0Sknfe+K8JtJlvgZH\nLkwDR3YJRvjiGDgNXnzPJKnTp1uL3ILNUVxQLpfTdoMni9t78TRG2okaRUIRKUnb\nOkUvKhQRq5MovKG95J9bk3PAAL4ZunQFR9WpZ8oXLkUUZn5UYx8wsWypoeQgdSnl\nDavppGGXGvQbbEehWNOxMQUAdC+01wiQbLCxuiL0IQKBgQDt9LSQ0yvUbqaP3aCC\njPfGVhXRn61oQRvwiY4xPkEI6kIKOry3NWx2TrNni7azH5kJuL39aAbbAs4ZEkZ0\nVcEAeMc8/A+tU9CJI9jYdNAoidKSF4umB5VsHEFWu562ZKFX7pkf0/oNryusO42N\nQr07xdkLj/zcqpTM1hkMauezcwKBgQDj+ia3zziF1FXskCKqzi0Vrbg4LEc95+ip\nKn6v8dB7eIIkcOZ/dJV96BooCPgMTcMQQZ9SXhTl9VX6410DJPdpeewrxNFFz6+R\nJXLA7w+s3xkFELe83pGCiLr4PpgsCLa2c/c40ZKmviQQWlYARrzz4ZwIMspVvWUx\n/4nHeNgnvwKBgQDbCCVxLcApgVS2FnYZ1XJ5UWIyG339+fii16tYLoqkjyUMz3ZI\nWdelqtm+1T7t8IbpCPFxIWE2FYXqVAUgLpNCQOp8ezRfAkhxA9slm5jDx+FY8m1r\n/Y0P/44xLPBxyX0uOIUgY/nGwNg9aG/qeGVBcONRnk9OO4JObkCOSpVOewKBgBrG\nWV/DHZ9SJdlHwFqRJmhiY71tLdzObNvZWtGtM9AcgiRnghq8BYapCNFk5EUFqQAy\nxNR2qjurybJSm6zv3YLsx9kIH4/0aWlBna2dJhkmhpi6sumitjI/fr4DE/ov982L\n5yAsNO4SbMfi5DDaHf0CpUGtPWg+ezEZEwRzY+efAoGBAO20C6ZWZgEDYZ0jA+zB\n8dRnHSpgOUUMszFeR6YaFXsQkag1UDA7CXJXSVECgehvzKNUQlmcMz3SgQ7l/tis\nusTm3CiEUIaOqOTTxRxbwdnwFNKgHNGbhFhdOITWg+AWf9DUHK4fLjF9VXCYxW1B\ncRamVY2tW7RRRKVwNyo71yYE\n-----END PRIVATE KEY-----"

type ITicketController interface {
	CreateTicket(c *fiber.Ctx) error
	GetTickets(c *fiber.Ctx) error
	GetTicketByID(c *fiber.Ctx) error
	//UpdateTicket(c *fiber.Ctx) error
	//DeleteTicket(c *fiber.Ctx) error
}

type ticketController struct {
	ticketRepo repos.ITicketRepository
}

func NewTicketController(ticketRepo repos.ITicketRepository) ITicketController {
	return &ticketController{ticketRepo}
}

func (t *ticketController) CreateTicket(c *fiber.Ctx) error {
	var ticket *models.Ticket

	// Create cancellable context.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	storageType := form.Value["storage_type"][0]

	r2storage, ociStorage, err := cloudStorage.GetStorageClient(storageType, accountId, r2BucketName, accessKeyId, accessKeySecret, tenancyOCID, userOCID, region, fingerprint, privateKey, namespace, ociBucketName)
	if err != nil {
		log.Print(err)
	}

	var attachments []*models.TicketAttachment

	attachmentFiles := form.File["attachments"]
	fmt.Println(attachmentFiles)
	for _, file := range attachmentFiles {
		fileName := file.Filename
		fileType := file.Header.Get("Content-Type")
		fileSize := file.Size

		// Create a buffer to read the file data.
		fileData, err := file.Open()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(err)
		}
		err = fileData.Close()
		if err != nil {
			return err
		}

		var fileUrl *string

		if storageType == "r2" {
			r2Client, err := r2storage.R2Init()
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON("r21")
			}
			fileUrl, err = r2storage.Upload(r2Client, fileName, fileType, fileData, c)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON("r22")
			}
		} else if storageType == "oci" {
			ociClient, err := ociStorage.OCIInit()
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON("oci1")
			}
			fileUrl, err = ociStorage.Upload(ctx, ociClient, fileName, fileType, fileSize, fileData)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON("oci2")
			}
		}

		attachment := &models.TicketAttachment{
			FileName: fileName,
			FileUrl:  *fileUrl,
			FileType: fileType,
			FileSize: int(fileSize),
		}

		attachments = append(attachments, attachment)
	}

	ticket = &models.Ticket{
		Subject:     form.Value["subject"][0],
		Description: form.Value["description"][0],
		Attachments: attachments,
	}

	valid, err := util.ValidateStruct(ticket)
	if !valid {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": govalidator.ErrorsByField(err)})
	}

	ticket, err = t.ticketRepo.CreateTicket(ctx, ticket)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "ticket has been created successfully!",
		"data":    ticket,
	})
}

func (t *ticketController) GetTickets(c *fiber.Ctx) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tickets, err := t.ticketRepo.GetTickets(ctx)
	if err != nil && err.Error() == "record not found" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"error":   err.Error(),
			"message": "tickets not found!",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"error":   err.Error(),
			"message": "something went wrong!",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "tickets have been fetched successfully!",
		"data":    tickets,
	})
}

func (t *ticketController) GetTicketByID(c *fiber.Ctx) error {

	// Create cancellable context.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Fetch ticket ID from request parameter
	id := c.Params("id")

	// Check if ticket ID is valid
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "please specify a valid ticket ID!",
		})
	}

	// Convert ticket ID to UUID
	idToUuid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "fail",
			"error":   err.Error(),
			"message": "please specify a valid ticket ID!",
		})
	}

	// Fetch ticket from database
	ticket, err := t.ticketRepo.GetTicketByID(ctx, idToUuid)
	if err != nil && err.Error() == "record not found" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "ticket not found!",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"error":   err.Error(),
			"message": "something went wrong!",
		})
	}

	return c.Status(fiber.StatusOK).JSON(ticket)
}

//func (t *ticketController) UpdateTicket(c *fiber.Ctx) error {
//
//	// Create cancellable context.
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	// Fetch ticket ID from request parameter
//	id := c.Params("id")
//
//	// Check if ticket ID is valid
//	if id == "" {
//		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
//			"status":  "fail",
//			"message": "Please specify a valid ticket ID!",
//		})
//	}
//
//	// Convert ticket ID to UUID
//	idToUuid, err := uuid.Parse(id)
//	if err != nil {
//		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
//			"status":  "fail",
//			"error":   err.Error(),
//			"message": "please specify a valid ticket ID!",
//		})
//	}
//
//	// Parse request body as a map of fields to update
//	var ticket *models.Ticket
//	if err := c.BodyParser(&ticket); err != nil {
//		return c.Status(http.StatusBadRequest).JSON(err)
//	}
//
//	updatedTicket, err := t.ticketRepo.UpdateTicket(ctx, idToUuid, ticket)
//	if err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(err)
//	}
//
//	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
//		"status":  "success",
//		"message": "Ticket has been updated successfully!",
//		"data":    updatedTicket,
//	})
//}

//
//func (t *ticketController) DeleteTicket(c *fiber.Ctx) error {
//	// Create cancellable context.
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	// Fetch url parameter.
//	id := c.Params("id")
//
//	if id == "" {
//		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
//			"status":  "fail",
//			"message": "Please specify a valid ticket ID!",
//		})
//	}
//
//	idToUuid, err := uuid.Parse(id)
//	if err != nil {
//		log.Println("Error parsing UUID:", err)
//		return nil
//	}
//
//	ticketID, err := t.ticketRepo.DeleteTicket(ctx, idToUuid)
//	if err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(err)
//	}
//
//	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
//		"status":  "success",
//		"message": "Ticket has been deleted successfully!",
//		"data":    ticketID,
//	})
//}
