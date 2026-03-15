package kinds

manifest: {
	appName:       "todo"
	groupOverride: "todo.grafana.app"

	versions: {
		"v1": v1
	}

	extraPermissions: {
		accessKinds: []
	}
}

v1: {
	kinds: [todov1]
	served: true
	codegen: {
		ts: enabled: true
		go: enabled: true
	}
}
