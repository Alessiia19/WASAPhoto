package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Login
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// User
	rt.router.PUT("/users/:userid", rt.wrap(rt.setMyUserName))
	rt.router.GET("/users/:userid", rt.wrap(rt.getUserProfile))
	rt.router.POST("/users/:userid/following", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:userid/following/:followingid", rt.wrap(rt.unfollowUser))
	rt.router.POST("/users/:userid/banned-users", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:userid/banned-users/:banneduserid", rt.wrap(rt.unbanUser))
	rt.router.GET("/users/:userid/banned-users/:banneduserid", rt.wrap(rt.getBanStatus))
	rt.router.GET("/users/:userid/stream", rt.wrap(rt.getMyStream))
	rt.router.GET("/users", rt.wrap(rt.getUsers))

	// Photo
	rt.router.POST("/users/:userid/photos", rt.wrap(rt.uploadPhoto))
	rt.router.POST("/users/:userid/photos/:photoid/likes", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/users/:userid/photos/:photoid/likes/:likeid", rt.wrap(rt.unlikePhoto))
	rt.router.POST("/users/:userid/photos/:photoid/comments", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/users/:userid/photos/:photoid/comments/:commentid", rt.wrap(rt.uncommentPhoto))
	rt.router.DELETE("/users/:userid/photos/:photoid", rt.wrap(rt.deletePhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
