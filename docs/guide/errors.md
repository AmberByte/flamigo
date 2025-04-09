# Error Handling

Flamigo follows a **structured and layered approach** to error handling. Errors are either **handled at their origin** or **passed down the stack**, wrapped with additional context at each level until they reach a layer that can appropriately handle or report them — often a strategy, service, or actor.

This approach makes it easy to trace where an error occurred while still enabling **clear, user-facing error messages**.

To support this pattern, Flamigo introduces its own error type: `flamigo.Error`.

:::tip  
The `flamigo.Error` type is fully compatible with Go’s built-in `error` interface — you can use it anywhere a standard `error` is expected.  
:::

---

## Creating a New Error

To define a new, custom error:

```go
flamigo.NewError(message string, opts ...ErrorOpt) *Error
```

This is ideal for crafting application-level or public-facing error messages.

---

## Wrapping an Existing Error

To enrich an existing error with additional context:

```go
flamigo.WrapError(message string, err error, opts ...ErrorOpt) *Error
```

This preserves the original error while adding a meaningful trace of where and why it failed — perfect for wrapping lower-level errors from databases, APIs, or other systems.

---

## Unwrapping Errors

Flamigo errors implement Go’s standard `Unwrap()` method, making them compatible with functions like `errors.Is()` and `errors.As()`:

```go
inner := errors.Unwrap(err)
```

You can traverse the error chain to inspect or extract specific layers, including public-facing error content.

---

## Public Error Messages

Flamigo errors are designed like **onions** — each layer adds more context.

At the core, a `flamigo.Error` may include a **public-facing message**, intended to be forwarded to clients. Higher-level functions (like strategies) can **wrap lower-level technical errors** and attach public responses using the `WithPublicResponse()` option.
By using the outermost PublicMessage() you get higher level erorr information without showing long internal error messages.

### Example

- A repository might return:  
  `"could not find entity"`
- A strategy wraps it as:  
  `"listing messages: could not find entity"`
- And attaches a public response:  
  `"Failed to list messages. Please try again later."`

Later, when unwrapping the error chain, this public response can be extracted and forwarded to the user — all while preserving the internal trace for logging or debugging.

This pattern enables granular error tracing **and** meaningful end-user communication.

---

## Error Message Formatting (Developer-Friendly)

Flamigo encourages using a **path-like, layered format** when writing internal error messages. This style helps developers trace errors through multiple layers of the application.

### ✅ Preferred:

```
strategy(foo:bar): listing messages: find(db_messages): context: deadline exceeded
```

### 🚫 Avoid:

```
failed strategy: could not list messages: failed finding in database: context: deadline exceeded
```

The preferred format enhances readability and helps you quickly pinpoint where and why something failed — especially useful when logs span multiple domains or systems.
