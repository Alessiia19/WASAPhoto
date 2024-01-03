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
	rt.router.PUT("/users/:userid", rt.wrap(rt.setMyUserName))
	rt.router.POST("/users/:userid/photos", rt.wrap(rt.uploadPhoto))
	rt.router.PUT("/users/:userid/following", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:userid/following/:followingid", rt.wrap(rt.unfollowUser))
	rt.router.PUT("/users/:userid/banned_users", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:userid/banned_users/:banneduserid", rt.wrap(rt.unbanUser))
	rt.router.GET("/users/:userid", rt.wrap(rt.getUserProfile))
	rt.router.GET("/users/:userid/stream", rt.wrap(rt.getMyStream))

	// Photo
	rt.router.PUT(("/users/:userid/photos/:photoid/likes"), rt.wrap(rt.likePhoto))
	rt.router.DELETE(("/users/:userid/photos/:photoid/likes/:likeid"), rt.wrap(rt.unlikePhoto))
	rt.router.POST(("/users/:userid/photos/:photoid/comments"), rt.wrap(rt.commentPhoto))
	rt.router.DELETE(("/users/:userid/photos/:photoid/comments/:commentid"), rt.wrap(rt.uncommentPhoto))
	rt.router.DELETE(("/users/:userid/photos/:photoid"), rt.wrap(rt.deletePhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
