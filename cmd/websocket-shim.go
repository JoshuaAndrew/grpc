package main

import (
	"os/exec"
	"fmt"
	"os"
	"io/ioutil"
	"log"
)
/*

cmd.Start 与 cmd.Wait 必须一起使用。
cmd.Start 不用等命令执行完成，就结束
cmd.Wait  等待命令结束
 */

func main() {

	args := os.Args
	var cmd *exec.Cmd
	if len(args) == 1 {
		fmt.Println("Usage:")
		return
	} else if len(args) == 2 {
		cmd = exec.Command(args[1])
	} else {
		cmd = exec.Command(args[1], args[2:]...)
	}
	//dateOut, err := cmd.Output()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(dateOut))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	stdout, err := cmd.Output()
	//stdout, err := cmd.StdoutPipe()  //指向cmd命令的stdout
	//defer stdout.Close()             //关闭输出流
	//// 运行命令
	//if err := cmd.Start(); err != nil {
	//	log.Fatal(err)
	//}



	// 读取输出结果
	content, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(content))

	//if err := cmd.Wait(); err != nil {
	//	log.Fatal(err)
	//}
}
