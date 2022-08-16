diff --git a/pkg/build/cmd.go b/pkg/build/cmd.go
index bd0128d..b4d7550 100644
--- a/pkg/build/cmd.go
+++ b/pkg/build/cmd.go
@@ -206,8 +206,8 @@ func ldflags(opts BuildOpts) (string, error) {

 	var b bytes.Buffer
 	b.WriteString("-w")
-	b.WriteString(fmt.Sprintf(" -X main.version=%s", opts.version))
-	b.WriteString(fmt.Sprintf(" -X main.commit=%s", getGitSha()))
+	b.WriteString(fmt.Sprintf(" -X main.version=%s", os.Getenv("GRAFANA_VERSION_HERE")))
+	b.WriteString(fmt.Sprintf(" -X main.commit=%s", "fix_heatmap,feat_extra_vars"))
 	b.WriteString(fmt.Sprintf(" -X main.buildstamp=%d", buildStamp))
 	b.WriteString(fmt.Sprintf(" -X main.buildBranch=%s", getGitBranch()))
 	if v := os.Getenv("LDFLAGS"); v != "" {
