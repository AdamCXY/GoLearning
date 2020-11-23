#pragma once
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include <stdio.h>
typedef unsigned char byte;
typedef unsigned word;
// s specifies the per-round shift amounts
unsigned S[64] = {
	7, 12, 17, 22,  7, 12, 17, 22,  7, 12, 17, 22,  7, 12, 17, 22,
	5,  9, 14, 20,  5,  9, 14, 20,  5,  9, 14, 20,  5,  9, 14, 20,
	4, 11, 16, 23,  4, 11, 16, 23,  4, 11, 16, 23,  4, 11, 16, 23,
	6, 10, 15, 21,  6, 10, 15, 21,  6, 10, 15, 21,  6, 10, 15, 21
};

// Use binary integer part of the sines of integers (Radians) as constants:
unsigned K[64] = {
	//0~15
	0xd76aa478, 0xe8c7b756, 0x242070db, 0xc1bdceee,
	0xf57c0faf, 0x4787c62a, 0xa8304613, 0xfd469501,
	0x698098d8, 0x8b44f7af, 0xffff5bb1, 0x895cd7be,
	0x6b901122, 0xfd987193, 0xa679438e, 0x49b40821,
	//16~31
	0xf61e2562, 0xc040b340, 0x265e5a51, 0xe9b6c7aa,
	0xd62f105d, 0x02441453, 0xd8a1e681, 0xe7d3fbc8,
	0x21e1cde6, 0xc33707d6, 0xf4d50d87, 0x455a14ed,
	0xa9e3e905, 0xfcefa3f8, 0x676f02d9, 0x8d2a4c8a,
	//32~47
	0xfffa3942, 0x8771f681, 0x6d9d6122, 0xfde5380c,
	0xa4beea44, 0x4bdecfa9, 0xf6bb4b60, 0xbebfbc70,
	0x289b7ec6, 0xeaa127fa, 0xd4ef3085, 0x04881d05,
	0xd9d4d039, 0xe6db99e5, 0x1fa27cf8, 0xc4ac5665,
	//48~63
	0xf4292244, 0x432aff97, 0xab9423a7, 0xfc93a039,
	0x655b59c3, 0x8f0ccc92, 0xffeff47d, 0x85845dd1,
	0x6fa87e4f, 0xfe2ce6e0, 0xa3014314, 0x4e0811a1,
	0xf7537e82, 0xbd3af235, 0x2ad7d2bb, 0xeb86d391
};



unsigned leftRotate(unsigned x, unsigned c) {
	return ((x << c) | (x >> (32 - c)));
}
//运算的时候按一个字节一个字节的来
void msgPadding(unsigned char* msg, int n, unsigned** words, int* lengthWords) {
	//计算填充后的msg的总长
	//至少添加一位，最多添加512位
	//最终的总长要为64的倍数（算上补上数据长度后的总长）
	//+1是因为取余运算的原因，这里我们必须向上取，所以需要多补一位
	//（+8）的原因是为了保证留下64位（8字节）添加数据长度（如果加上这8没有超过64的倍数则没影响，超过了就说明原来的长度也不够，需要多添一点，比如原长若只有60字节，则不能保证留下64位（8字节）空间给数据长度，这个+8使得原长溢出了64，说明了长度不够需要多补一位）
	int block = (n + 8) / 64 + 1;
	int finalLen = block * 64;

	//先添一个位的1，之后全0
	unsigned char* paddedMsg = (unsigned char*)malloc(finalLen);
	for (int i = 0; i < n; i++)
		paddedMsg[i] = msg[i];
	paddedMsg[n] = (unsigned char)0x80;
	for (int i = n + 1; i < finalLen - 8; i++)
		paddedMsg[i] = 0;

	//padding数据长度到最后两个字
	unsigned long long tmp = (unsigned long long)n * 8;

	memcpy(paddedMsg + finalLen - 8, &tmp, 8);

	int wordsLen = finalLen / 4;
	*lengthWords = wordsLen;
	*words = (unsigned*)malloc(finalLen);
	//printf("test1=%d\n", words);
	for (int i = 0; i < finalLen; i += 4)
	{
		unsigned thisWord = 0;
		for (int j = i; j <= i + 3; j++)
		{
			int shiftBits = (j - i) * 8;
			thisWord += ((unsigned)paddedMsg[j]) << shiftBits;
			
		}
		(*words)[i / 4] = thisWord;
	}
	//printf("test2=%d\n", *words);
	free(paddedMsg);
}

void ProcessMsg(unsigned* a, unsigned* b, unsigned* c, unsigned* d,unsigned* words, int lengthWords) {
	// Initialize variables:
	*a = 0x67452301;   // A
	*b = 0xefcdab89;   // B
	*c = 0x98badcfe;   // C
	*d = 0x10325476;   // D
	unsigned M[16];

	for (int i = 0; i < lengthWords / 16; i++) {
		unsigned A = *a;
		unsigned B = *b;
		unsigned C = *c;
		unsigned D = *d;
		for (int j = 0; j < 16; j++)
			M[j] = words[16 * i + j];

		for (int j = 0; j < 64; j++) {
			unsigned F, g;
			if (i >= 0 && i <= 15) {
				F = (B & C) | ((~B) & D);
				g = i;
			}
			else if (i >= 16 && i <= 31) {
				F = (D & B) | ((~D) & C);
				g = (5 * i + 1) % 16;
			}
			else if (i >= 32 && i <= 47) {
				F = B ^ C ^ D;
				g = (3 * i + 5) % 16;
			}
			else if (i >= 48 && i <= 63) {
				F = C ^ (B | (~D));
				g = (7 * i) % 16;
			}

			F = F + A + K[i] + M[g];
			A = D;
			D = C;
			C = B;
			B = B + leftRotate(F, S[i]);
		}
		a = a + A;
		b = b + B;
		c = c + C;
		d = d + D;
	}
}

unsigned char* MD5(unsigned char* msg) {
	unsigned * words;
	int wordsLen;
	msgPadding(msg, strlen(msg), &words, &wordsLen);
	printf("test3=%d\n", *words);
	
	unsigned  A, B, C, D;
	ProcessMsg(&A, &B, &C, &D, words, wordsLen);

	unsigned char* ans = (unsigned char*)malloc(33);
	ans[0] = '\0';

	unsigned char* a = (unsigned char*)&A, * b = (unsigned char*)&B,
		* c = (unsigned char*)&C, * d = (unsigned char*)&D;

	unsigned char* regs[] = { a,b,c,d };
	for (int i = 0; i < 4; i++)
	{
		unsigned char s[9];
		sprintf(s, "%02x%02x%02x%02x\0", regs[i][0], regs[i][1], regs[i][2], regs[i][3]);
		strcat(ans, s);
	}

	free(words);
	return ans;


}