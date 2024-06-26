package org.example.postplaceSpring.controller;

import org.example.postplaceSpring.service.CommentService;
import org.example.postplaceSpring.service.CustomUserDetails;
import org.example.postplaceSpring.service.PostService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.web.server.ResponseStatusException;

import java.io.IOException;
import java.util.UUID;

@RestController
public class CommentController {
    private final CommentService commentService;

    private static final Logger logger = LoggerFactory.getLogger(PostController.class);

    @Autowired
    public CommentController(CommentService commentService) {
        this.commentService = commentService;
    }

    @GetMapping("/public/getComments/{postId}")
    public ResponseEntity<String> getCommentsByPostId(@PathVariable UUID postId) {
        logger.info("Received Get request for /public/getComments/{commentId}");
        ResponseEntity<String> response = commentService.findCommentsByPostId(postId);
        if (response.getStatusCode().is2xxSuccessful()) {
            logger.info("Comment {} returned", postId);
            return response;
        } else {
            throw new ResponseStatusException(response.getStatusCode(), "Comment not found");
        }
    }

    @PostMapping("/user/comments/upload")
    public ResponseEntity<String> uploadComment( @RequestBody String commentJson) {
        logger.info("Received Post request for /user/comments/upload");
        // Get the authenticated user's details
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        CustomUserDetails userDetails = (CustomUserDetails) authentication.getPrincipal();
        long userId = userDetails.getUserId();
        // Pass the post file and userId to the service layer
        return commentService.createComment(commentJson, userId);
    }

    @DeleteMapping("/user/comments/delete/{commentId}")
    public ResponseEntity<Void> deleteComment(@PathVariable long commentId) {
        logger.info("Received Delete request for /user/comments/delete/{commentId}");
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        CustomUserDetails userDetails = (CustomUserDetails) authentication.getPrincipal();
        long userId = userDetails.getUserId();
        commentService.deleteCommentByCommentId(commentId, userId);
        return ResponseEntity.noContent().build();
    }

    @GetMapping("/public/getSubComments/{commentId}")
    public ResponseEntity<String> getSubCommentsByPostId(@PathVariable long commentId) {
        logger.info("Received Get request for /public/getSubComments/{commentId}");
        ResponseEntity<String> response = commentService.findSubCommentsByCommentId(commentId);
        if (response.getStatusCode().is2xxSuccessful()) {
            logger.info("Sub comment {} returned", commentId);
            return response;
        } else {
            throw new ResponseStatusException(response.getStatusCode(), "Sub comment not found");
        }
    }
    
    @PostMapping("/user/subComments/upload")
    public ResponseEntity<String> uploadSubComment( @RequestBody String subCommentJson) {
        logger.info("Received Post request for /user/subComments/upload");
        // Get the authenticated user's details
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        CustomUserDetails userDetails = (CustomUserDetails) authentication.getPrincipal();
        long userId = userDetails.getUserId();
        // Pass the post file and userId to the service layer
        return commentService.createSubComment(subCommentJson, userId);
    }

    @DeleteMapping("/user/subComments/delete/{subCommentId}")
    public ResponseEntity<Void> deleteSubComment(@PathVariable long subCommentId) {
        logger.info("Received Delete request for /user/subComments/delete/{subCommentId}");
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        CustomUserDetails userDetails = (CustomUserDetails) authentication.getPrincipal();
        long userId = userDetails.getUserId();
        commentService.deleteSubCommentBySubCommentId(subCommentId, userId);
        return ResponseEntity.noContent().build();
    }


}
