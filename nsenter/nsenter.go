package nsenter

/*
#include <errno.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>

__attribute__((constructor)) void enter_namespace(void) {
		char *zocker_pid;
		zocker_pid = getenv("zocker_pid");
		if (zocker_pid) {
			fprintf(stdout, "got zocker_pid=%s\n", zocker_pid);
		} else {
			// fprintf(stdout, "missing zocker_pid env skip nsenter\n");
			return;
		}
		char *zocker_cmd;
		zocker_cmd = getenv("zocker_cmd");
		if (zocker_cmd) {
			fprintf(stdout, "got zocker_cmd=%s\n", zocker_cmd);
		} else {
			fprintf(stdout, "missing zocker_cmd env skip nsenter\n");
			return;
		}
		int i;
		char nspath[1024];
		char *namespace[] = { "ipc", "uts", "net", "pid", "mnt" };
		for (i=0; i<5; i++) {
				sprintf(nspath, "/proc/%s/ns/%s", zocker_pid, namespace[i]);
				int fd = open(nspath, O_RDONLY);

				if (setns(fd, 0) == -1) {
					fprintf(stderr, "setns on %s namespace failed, %s\n", namespace[i], strerror(errno));
				} else {
					fprintf(stdout, "setns on %s namespace succeeded\n", namespace[i]);
				}
				close(fd);
		}
		int res = system(zocker_cmd);
		exit(0);
		return;
}
*/
import "C"
