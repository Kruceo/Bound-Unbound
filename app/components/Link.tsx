import React from "react";
import { usePageContext } from "vike-react/usePageContext";

export function Link(props: React.AnchorHTMLAttributes<HTMLAnchorElement>) {
  const { href, className, children, ...rest } = props
  const pageContext = usePageContext();
  const { urlPathname } = pageContext;

  const isActive = href === "/" ? urlPathname === href : urlPathname.startsWith(href??"");
  return (
    <a href={href} {...rest} className={className ? className : "" + (isActive ? "is-active" : undefined)}>
      {children}
    </a>
  );
}
