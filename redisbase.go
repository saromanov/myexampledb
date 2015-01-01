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

  	 m.Post("/", binding.Form(Item{}), func(item Item, ren render.Render) string{
  	 		newalbum := newAlbum(item.Band, item.Album, item.Members, item.Year)
  	 		r.Cmd("sadd", "testfield", newalbum)
  	 		return "FUN";
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