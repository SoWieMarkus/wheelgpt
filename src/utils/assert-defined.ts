/**
 * Asserts that a value is defined, throwing an error if it is null or undefined.
 * This is the preferred why over the `!` operator of typescript.
 *
 * @param value - The value to be checked.
 * @param name - The name of the value (optional). Defaults to "value".
 *
 * @throws {NotDefinedError} If the value is null or undefined.
 */
export function assertIsDefined<T>(
    value: T,
    name = "value",
): asserts value is NonNullable<T> {
    if (value === null || value === undefined) {
        throw new Error(`${name} is not defined.`);
    }
}