# Tutorial: Using Primary Keys in TSD

This hands-on tutorial walks you through using primary keys to build a blog management system with TSD.

---

## Table of Contents

- [Overview](#overview)
- [Scenario: Blog Management System](#scenario-blog-management-system)
- [Step 1: Define Types with Primary Keys](#step-1-define-types-with-primary-keys)
- [Step 2: Add Sample Data](#step-2-add-sample-data)
- [Step 3: Write Rules](#step-3-write-rules)
- [Step 4: Understand the Generated IDs](#step-4-understand-the-generated-ids)
- [Step 5: Query and Debug](#step-5-query-and-debug)
- [Step 6: Advanced Patterns](#step-6-advanced-patterns)
- [Best Practices Learned](#best-practices-learned)
- [Next Steps](#next-steps)

---

## Overview

By the end of this tutorial, you will:

✅ Understand how to define primary keys  
✅ Know when to use simple vs. composite keys  
✅ Use generated IDs in rules and conditions  
✅ Debug using IDs  
✅ Handle relationships between types  

**Time Required**: ~30 minutes  
**Difficulty**: Beginner

---

## Scenario: Blog Management System

We'll build a simple blog system with:

- **Users** - Authors and readers
- **Posts** - Blog articles
- **Comments** - User feedback on posts
- **Tags** - Categorization
- **PostTags** - Many-to-many relationship

---

## Step 1: Define Types with Primary Keys

### 1.1 User Type (Simple Primary Key)

Users have a natural unique identifier: their username.

```tsd
// Primary key: username (unique and stable)
type User(#username: string, email: string, role: string, created_at: number)
```

**Why `username` as primary key?**
- ✅ Naturally unique
- ✅ Stable (doesn't change often)
- ✅ Meaningful for debugging

**Generated ID format**: `User~alice`, `User~bob`

---

### 1.2 Post Type (Simple Primary Key)

Posts have a unique ID assigned by the system.

```tsd
// Primary key: post_id (unique identifier)
type Post(#post_id: string, author_username: string, title: string, content: string, published_at: number, status: string)
```

**Why `post_id` as primary key?**
- ✅ Guaranteed unique
- ✅ Explicitly identifies the post
- ✅ Better than composite (author + title might not be unique)

**Generated ID format**: `Post~POST-001`, `Post~POST-002`

---

### 1.3 Comment Type (Simple Primary Key)

Comments also have a unique ID.

```tsd
// Primary key: comment_id (unique identifier)
type Comment(#comment_id: string, post_id: string, author_username: string, text: string, created_at: number)
```

**Generated ID format**: `Comment~CMT-001`, `Comment~CMT-002`

---

### 1.4 Tag Type (Simple Primary Key)

Tags are identified by their name.

```tsd
// Primary key: name (unique tag name)
type Tag(#name: string, description: string, color: string)
```

**Generated ID format**: `Tag~programming`, `Tag~tutorial`

---

### 1.5 PostTag Type (Composite Primary Key)

PostTag is a junction table for the many-to-many relationship between Posts and Tags.

```tsd
// Composite primary key: post_id + tag_name
// A post can't have the same tag twice
type PostTag(#post_id: string, #tag_name: string, assigned_at: number)
```

**Why composite primary key?**
- ✅ Represents the relationship uniquely
- ✅ Prevents duplicate associations
- ✅ Natural combination of foreign keys

**Generated ID format**: `PostTag~POST-001_programming`, `PostTag~POST-001_tutorial`

---

### 1.6 View Type (No Primary Key - Hash)

Page views are temporal events without a natural key.

```tsd
// No primary key - uses hash
type View(post_id: string, viewer_username: string, viewed_at: number, duration: number)
```

**Why no primary key?**
- ❌ No natural unique identifier
- ✅ Same user can view same post multiple times
- ✅ Hash-based ID is fine for events

**Generated ID format**: `View~a1b2c3d4e5f6g7h8` (hash)

---

## Step 2: Add Sample Data

Create a file `blog_system.tsd`:

```tsd
// ============================================
// TYPE DEFINITIONS
// ============================================

type User(#username: string, email: string, role: string, created_at: number)
type Post(#post_id: string, author_username: string, title: string, content: string, published_at: number, status: string)
type Comment(#comment_id: string, post_id: string, author_username: string, text: string, created_at: number)
type Tag(#name: string, description: string, color: string)
type PostTag(#post_id: string, #tag_name: string, assigned_at: number)
type View(post_id: string, viewer_username: string, viewed_at: number, duration: number)

// ============================================
// SAMPLE DATA
// ============================================

// Users
User(username: "alice", email: "alice@blog.com", role: "author", created_at: 1700000000)
User(username: "bob", email: "bob@blog.com", role: "author", created_at: 1700001000)
User(username: "charlie", email: "charlie@blog.com", role: "reader", created_at: 1700002000)

// Posts
Post(post_id: "POST-001", author_username: "alice", title: "Getting Started with TSD", content: "TSD is a powerful rule engine...", published_at: 1700010000, status: "published")
Post(post_id: "POST-002", author_username: "alice", title: "Advanced RETE Algorithms", content: "The RETE algorithm...", published_at: 1700020000, status: "published")
Post(post_id: "POST-003", author_username: "bob", title: "Primary Keys Tutorial", content: "This tutorial explains...", published_at: 1700030000, status: "draft")

// Comments
Comment(comment_id: "CMT-001", post_id: "POST-001", author_username: "bob", text: "Great article!", created_at: 1700011000)
Comment(comment_id: "CMT-002", post_id: "POST-001", author_username: "charlie", text: "Very helpful, thanks!", created_at: 1700012000)
Comment(comment_id: "CMT-003", post_id: "POST-002", author_username: "charlie", text: "Can you explain more about beta nodes?", created_at: 1700021000)

// Tags
Tag(name: "tutorial", description: "Educational content", color: "blue")
Tag(name: "programming", description: "Code examples", color: "green")
Tag(name: "advanced", description: "Advanced topics", color: "red")

// Post-Tag relationships
PostTag(post_id: "POST-001", tag_name: "tutorial", assigned_at: 1700010100)
PostTag(post_id: "POST-001", tag_name: "programming", assigned_at: 1700010100)
PostTag(post_id: "POST-002", tag_name: "advanced", assigned_at: 1700020100)
PostTag(post_id: "POST-002", tag_name: "programming", assigned_at: 1700020100)
PostTag(post_id: "POST-003", tag_name: "tutorial", assigned_at: 1700030100)

// Views
View(post_id: "POST-001", viewer_username: "charlie", viewed_at: 1700015000, duration: 120)
View(post_id: "POST-001", viewer_username: "bob", viewed_at: 1700016000, duration: 90)
View(post_id: "POST-002", viewer_username: "charlie", viewed_at: 1700025000, duration: 180)
```

---

## Step 3: Write Rules

### 3.1 Find Published Posts by Author

```tsd
rule PublishedPostsByAuthor {
    when {
        u: User()
        p: Post()
        p.author_username == u.username
        p.status == "published"
    }
    then {
        print("Published post: '" + p.title + "' by " + u.username)
        print("  Post ID: " + p.id)
        print("  User ID: " + u.id)
    }
}
```

**Expected Output**:
```
Published post: 'Getting Started with TSD' by alice
  Post ID: Post~POST-001
  User ID: User~alice
Published post: 'Advanced RETE Algorithms' by alice
  Post ID: Post~POST-002
  User ID: User~alice
```

---

### 3.2 Find Posts with Comments

```tsd
rule PostsWithComments {
    when {
        p: Post()
        c: Comment()
        c.post_id == p.post_id
        p.status == "published"
    }
    then {
        print("Comment on '" + p.title + "': " + c.text)
        print("  Commenter: " + c.author_username)
    }
}
```

---

### 3.3 Find Popular Posts (Multiple Views)

```tsd
rule PopularPosts {
    when {
        p: Post()
        view_count: COUNT(v: View / v.post_id == p.post_id)
        view_count >= 2
    }
    then {
        print("Popular post: '" + p.title + "' with " + view_count + " views")
        print("  ID: " + p.id)
    }
}
```

---

### 3.4 Find Posts by Tag

```tsd
rule PostsByTag {
    when {
        t: Tag()
        pt: PostTag()
        p: Post()
        pt.tag_name == t.name
        pt.post_id == p.post_id
        t.name == "tutorial"
    }
    then {
        print("Tutorial post: '" + p.title + "'")
        print("  PostTag ID: " + pt.id)
    }
}
```

**Expected Output**:
```
Tutorial post: 'Getting Started with TSD'
  PostTag ID: PostTag~POST-001_tutorial
Tutorial post: 'Primary Keys Tutorial'
  PostTag ID: PostTag~POST-003_tutorial
```

---

## Step 4: Understand the Generated IDs

### ID Summary Table

| Type | Primary Key | Example ID |
|------|-------------|------------|
| User | `#username` | `User~alice` |
| Post | `#post_id` | `Post~POST-001` |
| Comment | `#comment_id` | `Comment~CMT-001` |
| Tag | `#name` | `Tag~tutorial` |
| PostTag | `#post_id, #tag_name` | `PostTag~POST-001_tutorial` |
| View | (none - hash) | `View~a1b2c3d4e5f6g7h8` |

### Why These ID Formats?

**Simple Keys** (`User~alice`):
- Easy to read and debug
- Predictable
- Natural uniqueness

**Composite Keys** (`PostTag~POST-001_tutorial`):
- Represents the relationship
- Prevents duplicates
- Contains both foreign keys

**Hash-Based** (`View~a1b2c3d4e5f6g7h8`):
- No natural key
- Deterministic (same data = same hash)
- Fine for temporal/event data

---

## Step 5: Query and Debug

### 5.1 Using IDs in Conditions

Find a specific post by ID:

```tsd
rule SpecificPost {
    when {
        p: Post()
        p.id == "Post~POST-001"
    }
    then {
        print("Found specific post: " + p.title)
    }
}
```

### 5.2 Debugging with IDs

Print all IDs for debugging:

```tsd
rule DebugAllPosts {
    when {
        p: Post()
    }
    then {
        print("Post: " + p.title)
        print("  ID: " + p.id)
        print("  Author: " + p.author_username)
        print("  Status: " + p.status)
    }
}
```

### 5.3 Tracing Relationships

Trace a comment back to its post and author:

```tsd
rule TraceComment {
    when {
        c: Comment()
        p: Post()
        u: User()
        c.post_id == p.post_id
        c.author_username == u.username
    }
    then {
        print("Comment chain:")
        print("  Comment ID: " + c.id)
        print("  Post ID: " + p.id)
        print("  Author ID: " + u.id)
    }
}
```

---

## Step 6: Advanced Patterns

### 6.1 Count Comments per Post

```tsd
rule PostCommentCount {
    when {
        p: Post()
        comment_count: COUNT(c: Comment / c.post_id == p.post_id)
        comment_count > 0
    }
    then {
        print("Post '" + p.title + "' has " + comment_count + " comments")
    }
}
```

### 6.2 Find Active Authors

Authors with published posts and comments:

```tsd
rule ActiveAuthors {
    when {
        u: User()
        u.role == "author"
        post_count: COUNT(p: Post / p.author_username == u.username AND p.status == "published")
        comment_count: COUNT(c: Comment / c.author_username == u.username)
        post_count > 0
        comment_count > 0
    }
    then {
        print("Active author: " + u.username)
        print("  Posts: " + post_count + ", Comments: " + comment_count)
        print("  User ID: " + u.id)
    }
}
```

### 6.3 Find Related Posts (Same Tags)

```tsd
rule RelatedPosts {
    when {
        pt1: PostTag()
        pt2: PostTag()
        p1: Post()
        p2: Post()
        pt1.tag_name == pt2.tag_name
        pt1.post_id == p1.post_id
        pt2.post_id == p2.post_id
        p1.post_id != p2.post_id
    }
    then {
        print("Related posts via tag '" + pt1.tag_name + "':")
        print("  Post 1: " + p1.title + " (ID: " + p1.id + ")")
        print("  Post 2: " + p2.title + " (ID: " + p2.id + ")")
    }
}
```

---

## Best Practices Learned

### ✅ DO

1. **Use natural keys when available**
   - `User(#username)` instead of `User(#id)`
   - Makes IDs readable and meaningful

2. **Use descriptive field names**
   - `#post_id` instead of `#id`
   - Clear what the field represents

3. **Use composite keys for relationships**
   - `PostTag(#post_id, #tag_name)`
   - Prevents duplicates naturally

4. **Use hash for events/logs**
   - `View` without primary key
   - Fine for temporal data

5. **Document your ID strategy**
   - Comment why each type uses its key
   - Helps future maintainers

### ❌ DON'T

1. **Don't use `id` as field name**
   - It's reserved for generated IDs
   - Use `post_id`, `user_id`, etc.

2. **Don't use unstable fields as keys**
   - Avoid `email` (might change)
   - Prefer `username` or system-generated IDs

3. **Don't use object/list types as keys**
   - Only primitive types allowed
   - string, number, bool only

4. **Don't set IDs explicitly**
   - TSD generates them automatically
   - Trust the system

---

## Next Steps

Now that you understand primary keys, try:

1. **Add more types**
   - `Category` with `#name`
   - `Like` with `#post_id, #username`
   - `Subscription` with `#subscriber, #author`

2. **Add more rules**
   - Find trending posts (most views in last hour)
   - Identify spam (many comments in short time)
   - Recommend posts (based on tags user views)

3. **Explore advanced features**
   - Aggregations with primary keys
   - Negation (`NOT`) with ID conditions
   - Retractions using IDs

4. **Read more documentation**
   - [Primary Keys Guide](../primary-keys.md)
   - [Architecture](../architecture/id-generation.md)
   - [API Reference](../api/id-generator.md)
   - [Migration Guide](../MIGRATION_IDS.md)

---

## Complete Code

The complete tutorial code is available in `examples/tutorials/blog_system.tsd`.

To run:

```bash
tsd compile examples/tutorials/blog_system.tsd
```

---

## Summary

In this tutorial, you learned:

✅ How to define primary keys (`#` prefix)  
✅ When to use simple vs. composite keys  
✅ When to use hash-based IDs (no primary key)  
✅ How to access IDs in rules (`fact.id`)  
✅ How IDs are formatted (`TypeName~identifier`)  
✅ How to debug using IDs  
✅ How to model relationships with composite keys  
✅ Best practices for choosing primary keys  

**Congratulations!** You now have a solid understanding of TSD's primary key system.

---

**Tutorial Version**: 1.0  
**Last Updated**: 2024-12-17  
**Difficulty**: Beginner  
**Time**: ~30 minutes