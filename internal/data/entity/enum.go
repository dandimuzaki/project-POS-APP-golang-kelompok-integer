package entity

// // UserRole enum
// type UserRole string

// const (
// 	RoleSuperAdmin UserRole = "superadmin"
// 	RoleAdmin      UserRole = "admin"
// 	RoleStaff      UserRole = "staff"
// )

// ProductStatus enum

// ProductAvailability enum
type ProductAvailability string

const (
	AvailabilityInStock    ProductAvailability = "in_stock"
	AvailabilityOutOfStock ProductAvailability = "out_of_stock"
)

// CustomerTitle enum
type CustomerTitle string

const (
	CustomerTitleMr   CustomerTitle = "Mr"
	CustomerTitleMrs  CustomerTitle = "Mrs"
	CustomerTitleMs   CustomerTitle = "Ms"
	CustomerTitleDr   CustomerTitle = "Dr"
	CustomerTitleProf CustomerTitle = "Prof"
)
