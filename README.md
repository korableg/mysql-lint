# mysql-lint
Check MySQL queries in your Golang code

__Steps for usage:__
1. Install linter
   ```bash
   go install github.com/korableg/mysql-lint@latest
   ```
2. Use it
   ```bash
   mysql-lint --dir my_awesome_go_project
   ```
   
If all your queries are OK, linter will exit with code 0, and you will see in your console something like that:
```
MySQL query linter: Start checking...
Checking finished successfully! 🎉
```
Otherwise, linter will exit with code 1 and will show all necessary information:
```
MySQL query linter: Start checking...
Found 2 errors:
test/database_sql/sql.go:79:32: issue: line 1 column 68 near "%s GROUP BY level ORDER BY level DESC LIMIT 15 AND K = 4" 
test/database_sql/sql.go:58:21: issue: line 1 column 49 near "%s)"
```

To skip checking the query you should add a __//nolint:mysql__ comment.

Now linter supports __database/sql__ interface only.

__Contributing:__  
We are glad to see everyone with pull requests :)
