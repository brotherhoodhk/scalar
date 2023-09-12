# BROTHERHOOD SCALAR TUTORIAL

### Basic Use
**Initialize Scalar Database**(make sure your gocache already opened)
```./scala init```
**start Scalar Main Server**
```./scala start```
### Configure
- **Server**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<scalar version="1.0">
    <gocache>
        <hostadd>localhost</hostadd>
        <port>8001</port>
        <default_db>scalar_db</default_db>
    </gocache>
    <plugins>
        <plugin_info>
        </plugin_info>
    </plugins>
    <paths>
        <common_path>/opt/scalar/plugins</common_path>
    </paths>
</scalar>
```
- **Client**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<scal-cli>
    <plugins>
        <plugin_info>
            <classname>reply parse</classname>
            <filename>replyparse.so</filename>
        </plugin_info>
        <plugin_info>
            <classname>reply parse</classname>
            <filename>zonename.so</filename>
        </plugin_info>
    </plugins>
    <common_path>/opt/scalar/plugins</common_path>
</scal-cli>
```