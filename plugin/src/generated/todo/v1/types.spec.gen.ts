// Code generated - EDITING IS FUTILE. DO NOT EDIT.

export interface Spec {
	// Title of the TODO item (required)
	title: string;
	// Optional description of the TODO item
	description?: string;
	// Current status of the TODO item
	status: "open" | "in_progress" | "done";
}

export const defaultSpec = (): Spec => ({
	title: "",
	status: "open",
});

