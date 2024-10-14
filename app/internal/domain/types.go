package domain

// Inline paths
const (
	InPathOrgs  = "/orgs"
	InPathRoles = "/roles"
)

// Web App paths
const (
	WebAppPathAdmin   = "/adminform"
	WebAppPathRegForm = "/regform"
	WebAppPathIntra   = "/intra"
	WebAppPathQrScan  = "/qr_scan"
)

// Web App (Front)
const (
	FormTypeGuestRegForm        = "type_guest_reg_form"
	FormTypeStudentRegForm      = "type_student_reg_form"
	FormTypeQR                  = "type_qr"
	FormTypeCreateAdmin         = "type_create_admin"
	FormTypeCreateInviteGuest   = "type_create_invite_guest"
	FormTypeCreateInviteStudent = "type_create_invite_student"
)

// Inline Query varchar(10)
const (
	InlEntranceGuest   = "g_entr"
	InlEntranceStudent = "s_entr"
)
