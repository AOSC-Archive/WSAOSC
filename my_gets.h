void my_gets_s(char* str, rsize_t n)
{
	int c = getchar();
	while (c != '\n' && c != EOF)
	{
		if (n > 1)
		{
			--n;
			*str++ = c;
		}

		c = getchar();
	}

	*str = '\0';
}
