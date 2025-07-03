CREATE UNIQUE INDEX addresses_unique_addr
    ON addresses (
        COALESCE(lower(trim(floor)),   ''),
        COALESCE(lower(trim(unit)),    ''),
        COALESCE(lower(trim(street)),  ''),
        number,
        COALESCE(lower(trim(city)),    ''),
        COALESCE(lower(trim(state)),   ''),
        COALESCE(trim(zip),            ''),
        upper(trim(country))
    );