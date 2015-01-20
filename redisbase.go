package main

import (
"fmt"
"github.com/fzzy/radix/redis"
"time"
"github.com/go-martini/martini"
"github.com/martini-contrib/render"
"github.com/codegangsta/martini-contrib/binding"
"strconv"
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
  	 		fmt.Println(item.Band)
  	 		fmt.Println(r.Cmd("type", item.Band))
  	 		r.Cmd("hmset", item.Band, newalbum)
  	 		r.Cmd("sadd", "title:band" + " " + item.Band + '"')
  	 		ren.HTML(200, "index", nil)
  	 	})

  	 m.Get("/find", func(ren render.Render){
  	 		ren.HTML(200, "find", nil)
  	 	})

  	 //Get list of all bands
  	 m.Get("/bands", func(ren render.Render){
  	 	results, err:= r.Cmd("hgetall", fnd.Search).Hash()
  	 	})

  	 m.Post("/find", binding.Form(Finder{}), func(fnd Finder, ren render.Render){
  	 		results, err:= r.Cmd("hgetall", fnd.Search).Hash()
  	 		value, _ := r.Cmd("smembers", "title:band" + fnd.Search).List()
  	 		fmt.Println("title:" + fnd.Search, value)
  	 		if err == nil {
  	 			ren.HTML(200, "find", results)
  	 		}
  	 		//resp, _ := results.List()
  	 		/*fmt.Println(resp)
  	 		fmt.Println(r.Cmd("scard", fnd.Search))
  	 		fmt.Println(r.Cmd("sinter", "Fun", fnd.Search))
  	 		//newmap := map[string]interface{}{"results": results}*/
  	 		//ren.HTML(200, "find", resp)
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

func newAlbum(band string, album string, members string, year int) map [string] string{
	return map[string]string {
		"Band": band,
		"Album": album,
		"Members": members /*strings.Split(members, ",")*/,
		"Year": strconv.Itoa(year),
	}
}
func getNames(item Item){
	fmt.Println(item.Band)
}
func main() {
    Init()
}