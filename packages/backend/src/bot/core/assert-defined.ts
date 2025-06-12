export function assertIsDefined<T>(value: T, name = "value"): asserts value is NonNullable<T> {
	if (value === null || value === undefined) {
		throw new Error(`${name} is not defined.`);
	}
}
