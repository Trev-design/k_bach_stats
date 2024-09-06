-module(pem_generator).

-include_lib("public_key/include/public_key.hrl").

-export([get_private_key/0]).

get_private_key() ->
    case file:read_file("private.pem") of
        {ok, Bin}  -> Bin;
        {error, _} -> generate_pems()
    end.


generate_pems() ->
    Priv = public_key:generate_key({rsa, 4096, 65537}),
    Pub = #'RSAPublicKey'{modulus = Priv#'RSAPrivateKey'.modulus, publicExponent = Priv#'RSAPrivateKey'.publicExponent},
    {_, PrivPem} = generate_pems_if_not_exist(Pub, Priv),
    PrivPem.


generate_pems_if_not_exist(Pub, Priv) ->
    PubPem = encoded_key_pem('RSAPublicKey', Pub),
    PrivPem = encoded_key_pem('RSAPrivateKey', Priv),

    generate_pem_file("public.pem", PubPem),
    generate_pem_file("private.pem", PrivPem),

    {PubPem, PrivPem}.


generate_pem_file(FilePath, Encoded) ->
    case file:read_file_info(FilePath) of
        {ok, _FileInfo} ->
            ok;
        _ -> 
            file:write_file(FilePath, Encoded)
    end.


encoded_key_pem(KeyType, RsaKey) ->
    Der = public_key:der_encode(KeyType, RsaKey),

    public_key:pem_encode([{'RSAPrivateKey', Der, not_encrypted}]).
