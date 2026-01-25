# Import Casdoor init resources

If you want to import all (to me known) init resources, run:

```shell
wget https://raw.githubusercontent.com/prochac/terraform-provider-casdoor/refs/heads/master/examples/import/resources_for_import.txt
wget https://raw.githubusercontent.com/prochac/terraform-provider-casdoor/refs/heads/master/examples/import/import_resources.sh
chmod +x import_resources.sh 
./import_resources.sh
```

And it will create `built-in.tf` HCL file with the resources.
