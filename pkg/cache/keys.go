package cache

import "fmt"

func CurrentDoctorAppointmentIDKey(doctorID uint) string {
	return fmt.Sprintf("doctor:%d:appointment_id", doctorID)
}

func AppointmentRoomIDKey(appointmentID string) string {
	return fmt.Sprintf("appointment:%s:room_id", appointmentID)
}

func RoomInfoKey(roomID string) string {
	return fmt.Sprintf("room:%s", roomID)
}
