# System Requirements - Lab CMS

## Overview

Lab CMS is a content management system designed for research laboratories. It provides a public-facing website for visitors to learn about the lab, its members, research, and publications, along with an administrative interface for lab members to manage content.

## User Roles

### Visitor
- Public users browsing the website without authentication
- Can view all public content (homepage, members, publications, news, projects)

### Normal Admin (Lab Member)
- Authenticated lab members with content management privileges
- Can create, edit, and delete their own content
- Cannot modify other admins' content or manage user accounts

### Root Admin
- Highest level administrative access
- Can create, edit, and delete any content (including other admins' content)
- Can add, remove, or edit other admin accounts
- Can assign or revoke admin privileges

---

## Public Website Requirements

### Homepage
- Display editable lab overview and introduction
- Show recent news/events
- Highlight featured research projects
- Provide navigation to other sections

### Lab Members
- Display list of current lab members organized by role:
  - **PI (Principal Investigator)** - Lab advisor/leader
  - **Postdocs** - Postdoctoral researchers
  - **PhD Students** - Doctoral students
  - **Master Students** - Master's degree students
  - **Bachelor Students** - Undergraduate students
  - **Researchers** - Other research staff
  - **Alumni** - Former members (optional section)
- Each member profile includes:
  - Name and role/category
  - Photo
  - Brief description/bio
  - Research interests (optional)
  - Link to personal page with more details
- Personal pages show comprehensive member information
- Members can be filtered or browsed by role

### Publications
- Display list of lab publications
- Simple listing format (title, authors, venue, year, link if available)
- Chronologically ordered (newest first)

### Research Projects
- Display ongoing and completed research projects
- Each project includes:
  - Title and description
  - Associated team members
  - Status (active/completed)
  - Relevant publications (linked)

### News & Events
- Display lab news, announcements, and events
- Chronologically ordered
- Include date, title, and content

---

## Admin System Requirements

### Authentication
- Login page for admin access
- Secure password-based authentication
- Session management

### Content Management Dashboard
- Overview of all managed content
- Quick access to create/edit functions
- Content status indicators

### Homepage Management
- Edit lab overview text
- Update featured content
- Manage homepage layout/sections

### Member Management
- Add new lab members with:
  - Name and role/category selection (PI, Postdoc, PhD, Master, Bachelor, Researcher)
  - Photos and bios
  - Research interests
  - Personal page content
- Edit existing member profiles
- Change member roles (e.g., when student becomes postdoc)
- Mark members as alumni when they leave
- Upload/manage member photos
- Create/edit personal pages
- Remove departed members or move to alumni section

### Publication Management
- Add new publications
- Edit publication details
- Remove outdated publications
- Link publications to projects/members

### Project Management
- Create new research projects
- Edit project descriptions and details
- Assign team members to projects
- Update project status
- Link related publications

### News Management
- Create news announcements
- Edit existing news
- Schedule or publish immediately
- Archive old news

### User Management (Root Admin Only)
- Add new admin accounts
- Remove admin accounts
- Edit admin permissions (normal vs root)
- Reset admin passwords

---

## User Stories

### Visitor Stories
- As a visitor, I want to see an overview of the lab on the homepage so I can understand what research the lab does
- As a visitor, I want to browse lab members so I can see who works there and their research interests
- As a visitor, I want to view detailed member profiles so I can learn more about specific researchers
- As a visitor, I want to see lab publications so I can read about their research output
- As a visitor, I want to browse research projects so I can understand the lab's current work
- As a visitor, I want to see recent news so I can stay updated on lab activities

### Normal Admin Stories
- As a lab member, I want to log in to the admin system so I can manage content
- As a lab member, I want to update my profile so visitors see current information
- As a lab member, I want to add my publications so they appear on the public site
- As a lab member, I want to edit content I created so I can fix mistakes or updates
- As a lab member, I want to create news posts so I can announce lab updates

### Root Admin Stories
- As a root admin, I want to manage the homepage content so it accurately represents the lab
- As a root admin, I want to add new lab members so the website stays current
- As a root admin, I want to edit any content so I can maintain quality across the site
- As a root admin, I want to add new admin accounts so new lab members can contribute
- As a root admin, I want to remove departed lab members so the website stays accurate
- As a root admin, I want to manage admin permissions so I can control access levels
