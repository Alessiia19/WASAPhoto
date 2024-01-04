package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Login
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// User
	rt.router.PUT("/users/:userID", rt.wrap(rt.setMyUserName))
	rt.router.POST("/users/:userID/photos", rt.wrap(rt.uploadPhoto))
	rt.router.POST("/users/:userID/following", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:userID/following/:followingID", rt.wrap(rt.unfollowUser))
	rt.router.POST("/users/:userID/banned_users", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:userID/banned_users/:bannedUserID", rt.wrap(rt.unbanUser))
	rt.router.GET("/users/:userID", rt.wrap(rt.getUserProfile))
	rt.router.GET("/users/:userID/stream", rt.wrap(rt.getMyStream))

	// Photo
	rt.router.POST(("/users/:userID/photos/:photoID/likes"), rt.wrap(rt.likePhoto))
	rt.router.DELETE(("/users/:userID/photos/:photoID/likes/:likeID"), rt.wrap(rt.unlikePhoto))
	rt.router.POST(("/users/:userID/photos/:photoID/comments"), rt.wrap(rt.commentPhoto))
	rt.router.DELETE(("/users/:userID/photos/:photoID/comments/:commentID"), rt.wrap(rt.uncommentPhoto))
	rt.router.DELETE(("/users/:userID/photos/:photoID"), rt.wrap(rt.deletePhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
