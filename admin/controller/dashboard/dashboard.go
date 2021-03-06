package dashboard

import(
	"net/http"
	"github.com/jschalkwijk/GolangBlog/admin/model/login"
	a "github.com/jschalkwijk/GolangBlog/admin/model/actions"
	"github.com/jschalkwijk/GolangBlog/admin/model/admin"
	"github.com/jschalkwijk/GolangBlog/admin/model/posts"
	"github.com/jschalkwijk/GolangBlog/admin/model/users"
	"github.com/jschalkwijk/GolangBlog/admin/controller"
)

func Index(w http.ResponseWriter, r *http.Request) {
	session := login.GetSession(r)

	if (!session.Logged) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	}

	if (r.PostFormValue("approve-selected") != ""){
		a.Approve(w,r,"posts")
	}
	if (r.PostFormValue("trash-selected") != ""){
		a.Trash(w,r,"posts")
	}
	if (r.PostFormValue("hide-selected") != ""){
		a.Hide(w,r,"posts")
	}
	data := new(admin.Data)
	data.Posts = *posts.All(0)
	data.Users = *users.All(0)
	data.Dashboard = true
	controller.View(w,"index", data)
}