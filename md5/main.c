#include<stdio.h>
//#include<cstring>
#include <string.h>
#include <ctype.h>
#include "md5.h"
unsigned char* MD5_cases[] = {
	"",
	"The quick brown fox jumps over the lazy dog",
	"a",
	"abc",
	"message digest",
	"abcdefghijklmnopqrstuvwxyz",
	"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789",
	"12345678901234567890123456789012345678901234567890123456789012345678901234567890"
};

int main() {
	//printf("%d", K[0]);
	unsigned int* words;
	int len;
	unsigned char* s[] = { " " };
	
	
	for (int i = 0; i < 8; i++) {
		printf("test4=%s\n",MD5(MD5_cases[i]));
	}
	
}