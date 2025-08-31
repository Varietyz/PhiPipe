// Copyright 2025 Baleine Jay
// Licensed under the Phicode Non-Commercial License (https://banes-lab.com/licensing)
// Commercial use requires a paid license. See link for details.
package server

import "github.com/gin-gonic/gin"

func (s *Server) Router() *gin.Engine {
	return s.router
}