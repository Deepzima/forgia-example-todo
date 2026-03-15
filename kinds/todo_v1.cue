package kinds

todov1: todoKind & {
	schema: {
		spec: {
			// Title of the TODO item (required)
			title: string
			// Optional description of the TODO item
			description?: string
			// Current status of the TODO item
			status: "open" | "in_progress" | "done"
		}
	}
}
