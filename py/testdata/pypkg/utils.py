def ml_run(s: str, b: bool, f: float, i: int) -> [str, bool, float, int]:
	print(f'ML got: s={s} b={b} f={f} i={i}')
	ns = s[::-1]
	nb = not b
	nf = 0 - f
	ni = 0 - i
	print(f'ML returned: s={ns} b={nb} f={nf} i={ni}')

	return ns, nb, nf, ni

def ml_image(bytes_data: bytes) -> float:
	if len(bytes_data) != 44080:
		raise Exception('not an image')

	return len(bytes_data)
