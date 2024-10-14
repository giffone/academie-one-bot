package domain

const (
	Started = "Start"
	Stopped = "Stop"
	Created = "Success!"

	AdminRightNeed  = "Administrator rights are required."
	NothingToDelete = "Not deleted."
	InviteIDWrong   = "Please check invite code. Can not find it."

	ErrorMsgClient = `Oops, something did not working properly.
Notify the administration of this error, please submit:
ChatID: %d
Date: %s`

	VerificationRegForm = `Success!
The administration needs time to verify your information. You will receive a notification.
After that you can scan the qr when you enter the academy campus and activate your session.`

	QRSuccessRegForm = `Success!
Welcome to campus.`

	QRSuccessStudent = `Success!
Now you can enter the Campus.
Your pass will expire at %s,
but you can renew it by scanning the QR again.`

	QRSuccessGuest = `Success!
Now you can enter the Campus.`

	IsBot = `Error!
Are you bot?`

	InviteForGuestAdmin = `Success!
You have created the guest invite code "%s" for "%s"!
The invitation will end on:
%s`
	InviteForStudentAdmin = `Success!
You have created the students invite code \"%s\"!
The invitation will end on:
%s.
You need to send out an invitation code to students in the mail and ask them to register as students on campus`

	InviteForGuestClient = `The invitation will end on: %s.
When visiting campus, you need to scan the QR code through this app by selecting "Campus Invitation" from the menu and clicking on the "Scan QR" button.`

	InviteExpired = `Your invitation has already expired on %s.
Please contact a member of staff %s.`

	RegNotConfirmed = `Your registration has not yet been confirmed.
Please contact a member of staff %s.`

	NoInvitation = `Oops!
Looks like you do not have an invitation.`

	MultipleInvitations = "You have %d invitations, choose one of them:"
)
