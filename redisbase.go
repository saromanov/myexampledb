package main

import (
"fmt"
"github.com/fzzy/radix/redis"
"time"
"github.com/go-martini/martini"
"github.com/martini-contrib/render"
"github.com/codegangsta/martini-contrib/binding"
)

//Fields for the base
type Item struct {
	Band string `form:"Band"`
	Album string `form:"Album"`
	Members string `form:"Members"`
	Year int `form:"Year"`
}

type Finder struct {
	Search string `form:"Finder"`
}

func Init(){
	r, err := InitRedis()
	fmt.Println(err)
	if err == nil{
		InitMartini(r)
	}
	closeRedis(r)
}
func InitMartini(r* redis.Client){
	 m := martini.Classic()
	 m.Use(render.Renderer())

  	 m.Get("/", func(ren render.Render){
    	ren.HTML(200,"index", nil)
  	 })

  	 m.Post("/", binding.Form(Item{}), func(item Item, ren render.Render){
  	 		newalbum := newAlbum(item.Band, item.Album, item.Members, item.Year)
  	 		r.Cmd("sadd", item.Band, newalbum)
  	 		fmt.Println("Album was append")
  	 		ren.HTML(200, "index", nil)
  	 	})

  	 m.Get("/find", func(ren render.Render){
  	 		ren.HTML(200, "find", nil)
  	 	})

  	 m.Post("/find", binding.Form(Finder{}), func(fnd Finder, ren render.Render){
  	 		results:= r.Cmd("smembers", fnd.Search);
  	 		resp, _ := results.List()
  	 		//newmap := map[string]interface{}{"results": results}
  	 		ren.HTML(200, "find", resp)
  	 	})

  	 m.NotFound(func(ren render.Render){
  	 		ren.HTML(200, "notfound", nil)
  	 	})
  	 m.Run()
}

func InitRedis() (*redis.Client, error){
	return redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
}

func closeRedis(r* redis.Client){
	defer r.Close()
}

func newAlbum(band string, album string, members string, year int) Item {
	return Item {
		Band: band,
		Album: album,
		Members: members /*strings.Split(members, ",")*/,
		Year: year,
	}
}
func getNames(item Item){
	fmt.Println(item.Band)
}
func main() {
    Init()
}