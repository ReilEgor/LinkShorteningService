package server

type Server interface {
	Run(port string) error
}
